package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"humn.ai/phan/mapbox"
	"humn.ai/phan/models"
	"humn.ai/phan/stdio"
	"humn.ai/phan/worker"
)

const (
	bufferSize          = 10000
	defaultWorkersCount = 5
)

// rootCmd represents the base command when called without any subcommands
var (
	workersCount int

	rootCmd = &cobra.Command{
		Use:   "app api_token [-w workers_count]",
		Short: "Read coordinates from stdin and print postcodes to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			// Token is passed to the first arg
			mbc := mapbox.NewClient(args[0])
			in := make(chan models.Input, bufferSize)
			out := make(chan models.Output, bufferSize)

			workerPool := worker.NewPool(mbc, in, out, workersCount)
			workerPool.Run()

			writeFinished := make(chan bool)
			go stdio.Write(os.Stdout, out, writeFinished)

			// Read will close the input channel when encounters EOF, which will eventually terminate workers
			go stdio.Read(os.Stdin, in)

			// Wait for workers
			workerPool.Wait()

			// Close output channel when all the workers finish, hence unblock Write goroutine
			close(out)
			<-writeFinished
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&workersCount, "workers-count", "w", defaultWorkersCount, "Set the number of workers")
}
