import type { ProjectData } from '@/api/project'
import { getProjects } from '@/api/project'
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarOpened = ref(true)
  const activeTabs = ref<any[]>([])
  const currentTab = ref('')

  // 项目相关状态
  const projectList = ref<ProjectData[]>([])
  const activeProjectAlias = ref(localStorage.getItem('active_project_alias') || '')

  function toggleSidebar() {
    sidebarOpened.value = !sidebarOpened.value
  }

  function addTab(tab: any) {
    if (!activeTabs.value.find(t => t.path === tab.path)) {
      activeTabs.value.push(tab)
    }
    currentTab.value = tab.path
  }

  function removeTab(path: string) {
    const index = activeTabs.value.findIndex(t => t.path === path)
    if (index !== -1) {
      activeTabs.value.splice(index, 1)
      if (currentTab.value === path && activeTabs.value.length > 0) {
        currentTab.value = activeTabs.value[activeTabs.value.length - 1].path
      }
    }
  }

  async function fetchProjects() {
    try {
      const data = await getProjects()
      projectList.value = data as unknown as ProjectData[]
      if (projectList.value.length > 0 && !activeProjectAlias.value) {
        activeProjectAlias.value = projectList.value[0].project_alias
      }
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    }
  }

  function setActiveProject(alias: string) {
    activeProjectAlias.value = alias
    localStorage.setItem('active_project_alias', alias)
  }

  // 监听 activeProjectAlias 变化，持久化到本地
  watch(activeProjectAlias, (newVal) => {
    if (newVal) {
      localStorage.setItem('active_project_alias', newVal)
    }
  })

  return {
    sidebarOpened,
    activeTabs,
    currentTab,
    projectList,
    activeProjectAlias,
    toggleSidebar,
    addTab,
    removeTab,
    fetchProjects,
    setActiveProject
  }
})

