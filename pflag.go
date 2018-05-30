package pflagsrc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/micro/go-config/source"
	"github.com/spf13/pflag"
)

var (
	defaultDelimiter = "-"
)

type flagsrc struct {
	delim string
	flags *pflag.FlagSet
	opts  source.Options
}

func (fs *flagsrc) Read() (*source.ChangeSet, error) {
	var changes map[string]interface{}
	fs.flags.VisitAll(func(f *pflag.Flag) {
		n := strings.ToLower(f.Name)
		keys := strings.Split(n, fs.delim)
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = f.Value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		mergo.Map(&changes, tmp) // need to sort error handling
		return
	})

	b, err := json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(b)
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	return &source.ChangeSet{
		Data:      b,
		Checksum:  checksum,
		Timestamp: time.Now(),
		Source:    fs.String(),
	}, nil
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func (fs *flagsrc) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

func (fs *flagsrc) String() string {
	return "flag"
}

// NewSource returns a config source that integrates the provided pflag.FlagSet.
// Hyphens are delimiters for nesting, and all keys are lowercased.
//
// Example:
//      dbhost := pflag.String("database-host", "localhost", "the db host name")
//
//      {
//          "database": {
//              "host": "localhost"
//          }
//      }
func NewSource(flags *pflag.FlagSet, opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	delim := defaultDelimiter

	if options.Context != nil {
		if d, ok := options.Context.Value(delimKey{}).(string); ok {
			delim = d
		}
	}

	return &flagsrc{flags: flags, opts: options, delim: delim}
}
