import { reactive } from 'vue'
import { GroupType, type Group } from '@/types/doris/analysis'

export function useGroupState(initialGroups: Group[] = []) {
  const groups = reactive<Group[]>(initialGroups)

  const addGroup = () => {
    groups.push({
      group_type: GroupType.Value,
      column: { table: 'event_data', field: '', alias: '' },
      value_ranges: [],
      time_grain: {
        column: { table: 'event_data', field: 'e_event_time', alias: '' },
        interval: 2,
        window_num: 0
      },
      tag_group: { tag_code: '', operator: 1, tag_value: '' },
      user_group: { group_name: '', group_code: '', operator: 1 }
    })
  }

  const removeGroup = (index: number) => {
    groups.splice(index, 1)
  }

  const clearGroups = () => {
    groups.splice(0, groups.length)
  }

  return {
    groups,
    addGroup,
    removeGroup,
    clearGroups
  }
}
