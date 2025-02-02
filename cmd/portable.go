package cmd

import (
	"path/filepath"

	"github.com/drakkan/sftpgo/dataprovider"
	"github.com/drakkan/sftpgo/service"
	"github.com/spf13/cobra"
)

var (
	directoryToServe             string
	portableSFTPDPort            int
	portableEnableSCP            bool
	portableAdvertiseService     bool
	portableAdvertiseCredentials bool
	portableUsername             string
	portablePassword             string
	portableLogFile              string
	portablePublicKeys           []string
	portablePermissions          []string
	portableCmd                  = &cobra.Command{
		Use:   "portable",
		Short: "Serve a single directory",
		Long: `To serve the current working directory with auto generated credentials simply use:

sftpgo portable

Please take a look at the usage below to customize the serving parameters`,
		Run: func(cmd *cobra.Command, args []string) {
			portableDir := directoryToServe
			if !filepath.IsAbs(portableDir) {
				portableDir, _ = filepath.Abs(portableDir)
			}
			service := service.Service{
				ConfigDir:     defaultConfigDir,
				ConfigFile:    defaultConfigName,
				LogFilePath:   portableLogFile,
				LogMaxSize:    defaultLogMaxSize,
				LogMaxBackups: defaultLogMaxBackup,
				LogMaxAge:     defaultLogMaxAge,
				LogCompress:   defaultLogCompress,
				LogVerbose:    defaultLogVerbose,
				Shutdown:      make(chan bool),
				PortableMode:  1,
				PortableUser: dataprovider.User{
					Username:    portableUsername,
					Password:    portablePassword,
					PublicKeys:  portablePublicKeys,
					Permissions: portablePermissions,
					HomeDir:     portableDir,
				},
			}
			if err := service.StartPortableMode(portableSFTPDPort, portableEnableSCP, portableAdvertiseService,
				portableAdvertiseCredentials); err == nil {
				service.Wait()
			}
		},
	}
)

func init() {
	portableCmd.Flags().StringVarP(&directoryToServe, "directory", "d", ".",
		"Path to the directory to serve. This can be an absolute path or a path relative to the current directory")
	portableCmd.Flags().IntVarP(&portableSFTPDPort, "sftpd-port", "s", 0, "0 means a random non privileged port")
	portableCmd.Flags().BoolVar(&portableEnableSCP, "scp", false, "Enable SCP")
	portableCmd.Flags().StringVarP(&portableUsername, "username", "u", "", "Leave empty to use an auto generated value")
	portableCmd.Flags().StringVarP(&portablePassword, "password", "p", "", "Leave empty to use an auto generated value")
	portableCmd.Flags().StringVarP(&portableLogFile, logFilePathFlag, "l", "", "Leave empty to disable logging")
	portableCmd.Flags().StringSliceVarP(&portablePublicKeys, "public-key", "k", []string{}, "")
	portableCmd.Flags().StringSliceVarP(&portablePermissions, "permissions", "g", []string{"list", "download"},
		"User's permissions. \"*\" means any permission")
	portableCmd.Flags().BoolVarP(&portableAdvertiseService, "advertise-service", "S", true,
		"Advertise SFTP/SCP service using multicast DNS")
	portableCmd.Flags().BoolVarP(&portableAdvertiseCredentials, "advertise-credentials", "C", false,
		"If the service is advertised via multicast DNS this flag allows to put username/password inside the advertised TXT record")
	rootCmd.AddCommand(portableCmd)
}
