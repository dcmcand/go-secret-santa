package csvparticipantloader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dcmcand/go-secret-santa/package/send"
)

type Loader struct{}

func (l *Loader) LoadParticipants(path string) (send.Participants, error) {
	file, err := os.Open(path)
	if err != nil {
		return send.Participants{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// Read header
	_, err = reader.Read()
	if err != nil {
		return send.Participants{}, fmt.Errorf("error reading header: %v", err)
	}
	p := send.Participants{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return send.Participants{}, fmt.Errorf("error reading file: %v\n", err)
		}
		participant := send.Participant{
			Name:      record[0],
			Email:     record[1],
			Partner:   record[2],
			Interests: strings.Split(record[3], ","),
		}
		p[record[0]] = participant
	}
	return p, nil

}
