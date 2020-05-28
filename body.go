package levi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Lv is the body of a leviathan request.
type Lv struct {
	echo.Context
	Logs []Log

	started  time.Time
	finished time.Time

	wg sync.WaitGroup
}

// Go performs an asynchronous job as part of a request.
//
// Request is not considered elapsed until all jobs are
// finished; request Logs are not getting flushed either.
func (lv *Lv) Go(job func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := errors.New(fmt.Sprintf("%v", r))
				lv.Panicf("%+v", err)
				lv.wg.Done()
			}
		}()

		job()
		lv.wg.Done()
	}()

	lv.wg.Add(1)
}

func (lv *Lv) Begin() {
	lv.started = time.Now()

	req := lv.Request()
	method, path := req.Method, req.URL.Path
	header := "INCOMING " + method + " " + path
	header += "NOW " + lv.started.Format(time.RFC3339)
	header += "ADDR " + lv.Addr()
	header += "AGENT " + lv.Agent()
	lv.inb4()
	header += "START"
	lv.log(PRINT, header)
}

func (lv *Lv) End(err error) {
	if err != nil {
		lv.Error(err)
	}

	lv.wg.Wait()
	lv.finished = time.Now()
	lv.logf(PRINT, "FINISHED WITH %d ELAPSED %s\n\n",
		lv.Response().Status,
		lv.finished.Sub(lv.started))
	logger.Report(lv)
}

// Addr reports the most likely Addr of the remote host.
func (lv *Lv) Addr() string {
	const cloudflare = "CF-Connecting-Addr"
	if ip, ok := lv.Request().Header[cloudflare]; ok {
		return ip[0]
	}

	return lv.Context.RealIP()
}

func (lv *Lv) Agent() string {
	return strings.Join(lv.Request().Header["User-Agent"], ", ")
}

// Putcookie sets an infinate root cookie.
func (lv *Lv) Putcookie(name, value string, httpOnly ...bool) {
	policy := http.SameSiteLaxMode
	domain := "localhost"

	if IsProd() {
		policy = http.SameSiteStrictMode
		domain = serverDomain
	}

	lv.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   IsProd(),
		HttpOnly: len(httpOnly) == 1 && httpOnly[0],
		SameSite: policy,
		Domain:   domain,
		Path:     "/",
	})
}

// Dropcookie removes a cookie.
func (lv *Lv) Dropcookie(name string) {
	lv.SetCookie(&http.Cookie{
		Name:     name,
		Secure:   IsProd(),
		HttpOnly: true,
		MaxAge:   -1,
	})
}

type Form interface {
	Validate(*Lv) error
	Apply(*Lv) error
}

func (lv *Lv) Paperwork(form Form) error {
	if err := lv.Bind(form); err != nil {
		lv.Error(err)
		return ErrBadPaperwork
	}

	if e := form.Validate(lv); e != nil {
		err, ok := e.(*ValidationError)
		if !ok {
			lv.Errorf("levi: when validating:", e)
		}

		return lv.JSON(http.StatusUnprocessableEntity, err)
	}

	return form.Apply(lv)
}

// Atomic runs a postgres transaction.
func (lv *Lv) Atomic(fn func(tx *pg.Tx) error) error {
	return db.RunInTransaction(fn)
}

// Table builds a new postgres orm query.
//
// Full postgres instance is usually not needed within
// the actual leviathan routes.
func (lv *Lv) Table(model interface{}) *orm.Query {
	return db.Model(model)
}

// Tables is like lv.Migrant(), but for slices.
func (lv *Lv) Tables(models ...interface{}) *orm.Query {
	return db.Model(models...)
}

// QueryInt64 works just like (*echo.Context).QueryInt, but with int64.
func (lv *Lv) QueryInt64(param string) (int64, error) {
	return strconv.ParseInt(lv.QueryParam(param), 10, 64)
}

func (lv *Lv) TBA() error {
	lv.Warn("TBA")
	return lv.NoContent(http.StatusNotImplemented)
}

func (lv *Lv) inb4() {
	if inb4 != nil {
		inb4(lv)
	}
}

// Kv is an arbitrary key-value data type, How can later
// be easily marshalled into JSON>
type Kv map[string]interface{}

// JSON returns the result of marshalling.
//
// By default, the output of this function is not pretty.
func (h Kv) Json(pretty ...bool) (string, error) {
	var (
		b   []byte
		err error
	)

	if len(pretty) != 0 && pretty[0] {
		b, err = json.MarshalIndent(h, "", "  ")
	} else {
		b, err = json.Marshal(h)
	}

	if err != nil {
		return "", fmt.Errorf("levi: failure to format : %w", err)
	}

	return string(b), nil
}
