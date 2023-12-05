package index

import (
	"testing"
)

func TestInsert(t *testing.T) {
	// bs1 := NewBackupSet("/backup/set/path1", "full")
	// bs2 := NewBackupSet("/backup/set/path2", "incr")
	// bs3 := NewBackupSet("/backup/set/path3", "incr")

	// bc := NewBackupCycle()

	// bc.Insert(bs1)
	// assert.Equal(t, 1, len(bc.BackupSets))

	// bc.Insert(bs2)
	// bc.Insert(bs3)
	// assert.Equal(t, 3, len(bc.BackupSets))

	// assert.Equal(t, "full", bc.Head().Type)
	// assert.Equal(t, "/backup/set/path1", bc.Head().Path)

	// assert.Equal(t, "incr", bc.Head().Next.Type)
	// assert.Equal(t, "/backup/set/path2", bc.Head().Next.Path)

	// assert.Equal(t, "incr", bc.Head().Next.Next.Type)
	// assert.Equal(t, "/backup/set/path3", bc.Head().Next.Next.Path)

	// jsonStr, err := bc.ToJSON()
	// if err != nil {
	// 	fmt.Println("Error converting BackupCycle to JSON:", err)
	// 	return
	// }

	// err = os.WriteFile("backup_cycle.json", []byte(jsonStr), 0644)
	// if err != nil {
	// 	fmt.Println("Error writing JSON to file:", err)
	// 	return
	// }

	// 从文件重新构建
	// fileData, err := os.ReadFile("backup_cycle.json")
	// if err != nil {
	// 	fmt.Println("Error reading JSON file:", err)
	// 	return
	// }

	// newBC, err := NewBackupCycleFromJSON(string(fileData))
	// if err != nil {
	// 	fmt.Println("Error creating BackupCycle from JSON:", err)
	// 	return
	// }

	// // 输出重新构建的 BackupCycle
	// fmt.Printf("Rebuilt BackupCycle:\n%+v\n", newBC)
}
