package server

import (
	"context"
	"os"

	pb "github.com/golimix/pm2-go/proto"
)

// start/update process
func (api *Handler) StartProcess(ctx context.Context, in *pb.StartProcessRequest) (*pb.Process, error) {
	api.mu.Lock()
	defer api.mu.Unlock()

	process := api.databaseById[in.Id]

	process.InitStartedAt()
	process.InitUptime()

	process.Name = in.Name
	process.ExecutablePath = in.ExecutablePath
	process.Args = in.Args
	process.PidFilePath = in.PidFilePath
	process.LogFilePath = in.LogFilePath
	process.ErrFilePath = in.ErrFilePath
	process.AutoRestart = in.AutoRestart
	process.Cwd = in.Cwd
	process.Pid = in.Pid
	process.CronRestart = in.CronRestart
	process.ProcStatus.ParentPid = 1
	err := process.UpdateNextStartAt()
	if err != nil {
		return nil, err
	}

	found, err := os.FindProcess(int(in.Pid))
	if err != nil {
		process.Pid = in.Pid
	}

	process.SetStopSignal(false)
	process.SetStatus("online")
	updateProcessMap(api, in.Id, found)

	return process, nil
}
