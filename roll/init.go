package roll

import "github.com/flywingedai/dice/core"

func Initialize() {

	// Creation
	core.AddRollType(core.ROLL_SIDES, func() core.Roll { return &roll_Sides{} })
	core.AddRollType(core.ROLL_WEIGHTED, func() core.Roll { return &roll_Weighted{} })

	// Rolls
	core.AddRollType(core.ROLL_MULTIPLE, func() core.Roll { return &roll_Multiple{} })

	// Skip for Merge
	core.AddRollType(core.ROLL_SKIP, func() core.Roll { return &roll_Skip{} })
}
