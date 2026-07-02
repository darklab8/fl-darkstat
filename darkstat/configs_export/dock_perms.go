package configs_export

type DockPerms int8

const (
	DockUnknown DockPerms = iota
	DockIFF
	DockIFFPlus
	DockIFFMinus
	DockACL
	DockACLPlus
	DockACLMinus
	DockNoAccess
)

func (d DockPerms) ToStr() string {
	switch d {
	case DockIFF:
		return "iff"
	case DockIFFPlus:
		return "iff+"
	case DockIFFMinus:
		return "iff-"
	case DockACL:
		return "acl"
	case DockACLPlus:
		return "acl+"
	case DockACLMinus:
		return "acl-"
	case DockNoAccess:
		return "unk"
	}
	return "nil"
}

func (d DockPerms) TypeToColor() string {
	switch d {
	case DockIFF:
		return "rgba(0, 219, 0, 0.3)"
	case DockIFFPlus:
		return "rgba(0, 255, 0, 0.3)"
	case DockIFFMinus:
		return "rgba(145, 177, 2, 0.3)"
	case DockACL:
		return "rgba(255, 153, 0, 0.3)"
	case DockACLPlus:
		return "rgba(251, 255, 0, 0.3)"
	case DockACLMinus:
		return "rgba(255, 94, 0, 0.3)"
	case DockNoAccess:
		return "rgba(255, 0, 0, 0.3)"
	}
	return "rgba(225, 0, 255, 0.3)"
}

/*
see this type of value for up to date info

	func (d DefenseMode) ToStr() string {
		switch d {

		case 1:
			return "SRP Whitelist > Blacklist > IFF Standing, Anyone with good standing"
		case 2:
			return "Whitelist > Nodock, Whitelisted ships only"
		case 3:
			return "Whitelist > Hostile, Whitelisted ships only"
		default:
			return "not recognized"
		}
	}
*/
func DockingPermissions(goods ...*MarketGood) DockPerms {
	dock_access := DockUnknown

	for _, good := range goods {
		if good.PoB == nil {
			if dock_access < DockIFF {
				dock_access = DockIFF
			}
		}

		if good.PoB != nil {

			if dock_result := DockingPermissionsPoB(&good.PoB.PoBCore); dock_access < dock_result {
				dock_access = dock_result
			}
		}
	}
	return dock_access
}

func DockingPermissionsPoB(pob *PoBCore) DockPerms {
	dock_access := DockUnknown

	if pob.DefenseMode == nil {
		return DockNoAccess
	}

	if *pob.DefenseMode == 2 || *pob.DefenseMode == 3 {
		if dock_access < DockACL {
			dock_access = DockACL

			if pob.HasDockFriendlyFactions && dock_access < DockACLPlus {
				dock_access = DockACLPlus
			}
			if pob.HasDockEnemyFactions && dock_access < DockACLMinus {
				dock_access = DockACLMinus
			}
		}
	} else if *pob.DefenseMode == 1 {
		if dock_access < DockIFF {
			dock_access = DockIFF

			if pob.HasDockFriendlyFactions && dock_access < DockIFFPlus {
				dock_access = DockIFFPlus
			}
			if pob.HasDockEnemyFactions && dock_access < DockIFFMinus {
				dock_access = DockIFFMinus
			}
		}
	}
	return dock_access
}
