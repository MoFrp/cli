package banner

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/version"
)

var (
	white  = lipgloss.Color("255")
	pink   = lipgloss.Color("201")
	whiteT = lipgloss.NewStyle().Foreground(white).Italic(true)
	pinkT  = lipgloss.NewStyle().Foreground(pink).Italic(true)

	cliBanner = `
      __  __         _____                ____  _      ___ 
     |  \/  |  ___  |  ___|_ __  _ __    / ___|| |    |_ _|
     | |\/| | / _ \ | |_  | '__|| '_ \  | |    | |     | | 
     | |  | || (_) ||  _| | |   | |_) | | |___ | |___  | | 
     |_|  |_| \___/ |_|   |_|   | .__/   \____||_____||___|
                                |_|                        
`

	nodeBanner = `
      __  __         _____                _   _             _           ____  _      ___ 
     |  \/  |  ___  |  ___|_ __  _ __    | \ | |  ___    __| |  ___    / ___|| |    |_ _|
     | |\/| | / _ \ | |_  | '__|| '_ \   |  \| | / _ \  / _' | / _ \  | |    | |     | | 
     | |  | || (_) ||  _| | |   | |_) |  | |\  || (_) || (_| ||  __/  | |___ | |___  | | 
     |_|  |_| \___/ |_|   |_|   | .__/   |_| \_| \___/  \__,_| \___|   \____||_____||___|
                                |_|                                                      
`
)

func DisplayCLIBanner() {
	fmt.Print(whiteT.Render(cliBanner))
	fmt.Println()
	log.Infof("Nya! MoFrp CLI %s 启动中...", version.Version())
}

func DisplayNodeBanner() {
	fmt.Print(pinkT.Render(nodeBanner))
	fmt.Println()
	log.Infof("Nya! MoFrp Node %s 启动中...", version.Version())
}
