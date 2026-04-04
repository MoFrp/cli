package banner

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/version"
)

var (
	white  = lipgloss.Color("255")
	whiteT = lipgloss.NewStyle().Foreground(white).Italic(true)

	cliBanner = `
      __  __         _____                ____  _      ___ 
     |  \/  |  ___  |  ___|_ __  _ __    / ___|| |    |_ _|
     | |\/| | / _ \ | |_  | '__|| '_ \  | |    | |     | | 
     | |  | || (_) ||  _| | |   | |_) | | |___ | |___  | | 
     |_|  |_| \___/ |_|   |_|   | .__/   \____||_____||___|
                                |_|                        
`
)

func DisplayCLIBanner() {
	fmt.Print(whiteT.Render(cliBanner))
	fmt.Println()
	log.Infof("Nya! MoFrp CLI %s 启动中...", version.Version())
}
