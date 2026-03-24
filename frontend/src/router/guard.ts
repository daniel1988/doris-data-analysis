import router from './index'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const appStore = useAppStore()

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - DMP Admin V2`
  }

  // 检查是否需要登录
  const token = userStore.token
  if (!to.meta.public && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    // 路由记录到 tabs
    if (to.path !== '/login') {
      appStore.addTab({
        path: to.path,
        title: to.meta.title as string
      })
    }
    next()
  }
})

export default router
