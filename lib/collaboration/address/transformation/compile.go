package transformation

import (
	"fmt"
	"os"

	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func ExecuteSqlTemplate(template string, pongoCtx pongo2.Context) (sqlString string, err error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		err = fmt.Errorf("generating sql template: %w", err)
		log.Error(err)
		return sqlString, err
	}

	sqlString, err = tpl.Execute(pongoCtx)
	if err != nil {
		err = fmt.Errorf("executing sql template: %w", err)
		log.Error(err)
		return sqlString, err
	}
	return sqlString, nil
}
