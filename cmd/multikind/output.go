package multikind

import (
	"encoding/csv"
	"errors"
	"io"

	"github.com/olekukonko/tablewriter"
)

type Format string

const (
	UnknownFormat Format = "unknown"
	Table         Format = "table"
	CSV           Format = "csv"
)

func MustParseFormat(s string) Format {
	switch s {
	case string(Table):
		return Table
	case string(CSV):
		return CSV
	default:
		return UnknownFormat
	}
}

func NewFormatWriter(w io.Writer, format Format) *FormatWriter {
	return &FormatWriter{
		w:      w,
		format: format,
	}
}

type FormatWriter struct {
	w      io.Writer
	format Format
}

func (f *FormatWriter) WriteAndClose(headers []string, items [][]string) error {
	if f.format == Table {
		table := tablewriter.NewWriter(f.w)
		table.SetHeader(headers)
		for _, item := range items {
			table.Append(item)
		}
		table.Render()
		return nil
	} else if f.format == CSV {
		ww := csv.NewWriter(f.w)
		ww.Write(headers)
		for _, item := range items {
			ww.Write(item)
		}
		ww.Flush()
		return ww.Error()
	} else {
		return errors.New("not implemented")
	}
}
