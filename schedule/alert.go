package schedule

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	htmltpl "html/template"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	texttpl "text/template"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/store"
	"github.com/jordan-wright/email"
)

type Alerter struct {
	ts       store.TaskStore
	js       store.JobStore
	cs       store.ConfigStore
	us       store.UserStore
	channels map[string]AlertChannel
}

func NewAlerter(ts store.TaskStore, js store.JobStore, cs store.ConfigStore, us store.UserStore) *Alerter {
	return &Alerter{
		ts: ts,
		js: js,
		cs: cs,
		us: us,
		channels: map[string]AlertChannel{
			"email": EmailChannel{},
			"wecom": WeComChannel{},
		},
	}
}

func (a *Alerter) Alert(jobId string, info string) {
	err := a.alert(jobId, info)
	if err != nil {
		log.Get("schedule").Errorf("alert failed: %s", err)
	}
}

func (a *Alerter) alert(jobId string, info string) error {
	job, err := a.js.Find(jobId)
	if err != nil {
		return err
	}

	task, err := a.ts.Find(job.Task)
	if err != nil {
		return err
	}

	if len(task.Alerts) == 0 {
		return nil
	}

	vars := data.Map{
		"error":     info,
		"task":      task.Name,
		"handler":   task.Handler,
		"runner":    task.Runner,
		"job":       job.Id.Hex(),
		"mode":      job.Mode,
		"fire":      job.FireTime.Format("2006-01-02 15:04:05"),
		"args":      job.Args,
		"scheduler": job.Scheduler,
	}
	if job.Execute.EndTime == nil {
		vars.Set("duration", "")
	} else {
		vars.Set("duration", time.Time(*job.Execute.EndTime).Sub(time.Time(*job.Execute.StartTime)).String())
	}

	users, err := a.us.Fetch(task.Maintainers)
	if err != nil {
		return err
	}

	for _, alert := range task.Alerts {
		ch := a.channels[alert]
		if ch == nil {
			log.Get("schedule").Warnf("unknown alert method: %s", alert)
			continue
		}

		options, err := a.cs.Find("alert." + alert)
		if err != nil {
			log.Get("schedule").Errorf("failed to fetch alert.%s options: %s", alert, err)
		}

		err = ch.Send(options, users, vars)
		if err != nil {
			log.Get("schedule").Errorf("failed to send %s alert: %s", alert, err)
		}
	}
	return nil
}

func transform(html bool, tpl string, data interface{}) (s string, err error) {
	var (
		t interface {
			Execute(wr io.Writer, data interface{}) error
		}
	)
	if html {
		t, err = htmltpl.New("").Parse(tpl)
	} else {
		t, err = texttpl.New("").Parse(tpl)
	}
	if err == nil {
		buf := &bytes.Buffer{}
		if err = t.Execute(buf, data); err == nil {
			s = buf.String()
		}
	}
	return
}

type AlertChannel interface {
	Send(options data.Options, users []*store.User, vars data.Map) (err error)
}

type EmailChannel struct {
}

func (c EmailChannel) Send(options data.Options, users []*store.User, vars data.Map) (err error) {
	const (
		defaultTitle = "[Skynet]Failed to execute task: {{ .task }}"
		defaultBody  = "Task: {{ .task }}，Error: {{ .error }}"
	)

	if options.Get("enabled") != "true" {
		return nil
	}

	var (
		sender       = options.Get("sender")
		receiver     = options.Get("receiver")
		smtpAddress  = options.Get("smtp_address")
		smtpUsername = options.Get("smtp_username")
		smtpPassword = options.Get("smtp_password")
		title        = options.Get("title")
		body         = options.Get("body")
	)

	var emails []string
	for _, user := range users {
		if user.Email != "" {
			emails = append(emails, user.Email)
		}
	}
	if len(emails) == 0 {
		if receiver == "" {
			return nil
		}
		emails = strings.Split(receiver, ",")
	}

	if smtpAddress == "" || sender == "" {
		return errors.New("missing options of SMTP")
	}
	if title == "" {
		title = defaultTitle
	}
	if body == "" {
		body = defaultBody
	}
	if title, err = transform(true, title, vars); err != nil {
		return err
	}
	if body, err = transform(true, body, vars); err != nil {
		return err
	}

	mail := &email.Email{
		To:      emails,
		From:    sender,
		Subject: title,
		HTML:    []byte(body),
		//Text:    []byte(body),
	}
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, "")
	return mail.SendWithTLS(smtpAddress, auth, &tls.Config{InsecureSkipVerify: true})
}

type WeComChannel struct {
}

func (c WeComChannel) Send(options data.Options, users []*store.User, vars data.Map) (err error) {
	const (
		defaultTitle = "[Skynet]Failed to execute task: {{ .task }}"
		defaultBody  = "Task: {{ .task }}，Error: {{ .error }}"
	)

	if options.Get("enabled") != "true" {
		return nil
	}

	var ids []string
	var maintainers []string
	for _, user := range users {
		if user.Wecom != "" {
			ids = append(ids, user.Wecom)
		}
		maintainers = append(maintainers, "@"+user.Name)
	}
	vars.Set("maintainers", strings.Join(maintainers, " "))

	var (
		mode    = options.Get("mode")
		msgType = options.Get("msg_type")
		title   = options.Get("title")
		body    = options.Get("body")
	)

	if title == "" {
		title = defaultTitle
	}
	if body == "" {
		body = defaultBody
	}
	if title, err = transform(false, title, vars); err != nil {
		return err
	}
	if body, err = transform(false, body, vars); err != nil {
		return err
	}

	args := data.Map{
		"msgtype": msgType,
		msgType: data.Map{
			"content":        body,
			"mentioned_list": ids,
		},
	}

	switch mode {
	case "robot":
		return c.sendRobot(options, args)
	case "app":
		return c.sendApp(options, args)
	case "":
		return errors.New("missing mode option")
	default:
		return errors.New("unknown mode: " + mode)
	}
}

func (c WeComChannel) sendRobot(options data.Options, args data.Map) (err error) {
	const robotUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

	robotKey := options.Get("robot_key")
	if robotKey == "" {
		return errors.New("missing robot_key option")
	}
	return c.post(robotUrl+robotKey, args)
}

func (c WeComChannel) sendApp(options data.Options, args data.Map) (err error) {
	return errors.NotSupported
}

func (c WeComChannel) post(url string, args data.Map) error {
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, web.MIMEApplicationJSON, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// response: {"errcode":0,"errmsg":"ok"}
	result := struct {
		Code int32  `json:"errcode"`
		Msg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	} else if result.Code != 0 {
		return errors.Coded(result.Code, result.Msg)
	}
	return nil
}
