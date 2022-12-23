package analytics

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"runtime/pprof"
)

func WriteMemoryProfile(filename string) error {
	log.Info().Msgf("creating memory profile: %s", filename)
	fmem, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create memory profile file: %w", err)
	}
	defer fmem.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		return fmt.Errorf("could not file memory profile: %w", err)
	}
	log.Info().Msgf("memory profile saved: %s", filename)
	return nil
}
