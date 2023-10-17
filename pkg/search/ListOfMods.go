package search

import (
	"fmt"
)

func ListOfMods(args []string) {
	modList := Search(args)

	for i, mod := range modList.Data{
			if mod.Category.Name != "" {
				mod.Category.Name = "/" + " " + mod.Category.Name
			}

			fmt.Printf(
				`
		[%v] %s 
			🧍 %s
			📍 %s %s
			❤  %v 😋 %v 👁️  %v		🕑 %s
		
		`, i+1,
				mod.Name,
				mod.User.Name,

				mod.Game.Name, 
				mod.Category.Name,

				mod.Likes,
				mod.Downloads,
				mod.Views,

				// parse the mod.UpdatedAt time relative like 10 days ago
				mod.UpdatedAt.Local().Format("Jan 2, 2006"),
			)
		}
}