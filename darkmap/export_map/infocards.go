package export_map

func (e *Export) GetFullInfocardIds(infocard_id int) []int {
	var infocard_ids []int = make([]int, 0)
	infocard_ids = append(infocard_ids, infocard_id)
	if infocard_middle_id, exists := e.Mapped.InfocardmapINI.InfocardMapTable.Map[infocard_id]; exists {
		infocard_ids = append(infocard_ids, infocard_middle_id)
	}
	return infocard_ids
}
