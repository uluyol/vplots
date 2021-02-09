package conversions

import (
	"flag"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

var debugExec = flag.Bool("dexec", false, "debug executed commands")

func loggingExec(logAttr, name string, args ...string) ([]byte, error) {
	if *debugExec {
		var buf strings.Builder
		buf.WriteString("exec: ")
		if logAttr != "" {
			buf.WriteString(logAttr)
			buf.WriteString(": ")
		}
		buf.WriteString(name)
		for _, a := range args {
			buf.WriteByte(' ')
			buf.WriteString(a)
		}
		log.Print(buf.String())
	}
	return exec.Command(name, args...).Output()
}

func externalToSVG(pdfPath string) ([]byte, error) {
	return loggingExec("", "pdf2svg", pdfPath, "/dev/stdout")
}

func (c *Converter) externalToPNG(pdfPath string, width int, logAttr string) ([]byte, error) {
	if width == 0 {
		return loggingExec(logAttr, "convert", "-density", "300", pdfPath, "png:-")
	}

	dim := c.getDim(pdfPath)
	dpiToUse := float64(width) / dim.wInches

	return loggingExec(logAttr, "convert", "-density", strconv.Itoa(int(math.Ceil(dpiToUse))), "-resize", strconv.Itoa(width), pdfPath, "png:-")
}
