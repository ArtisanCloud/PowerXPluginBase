<template>
  <div class=\"min-h-screen bg-gray-50\">
    <!-- 页面头部 -->
    <header class=\"bg-white shadow-sm border-b\">
      <div class=\"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8\">
        <div class=\"flex justify-between items-center py-6\">
          <div class=\"flex items-center\">
            <h1 class=\"text-2xl font-bold text-gray-900\">
              PowerX Scrum Plugin
            </h1>
          </div>
          <nav class=\"flex space-x-8\">
            <NuxtLink to=\"/\" class=\"text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium\">
              仪表板
            </NuxtLink>
            <NuxtLink to=\"/tasks\" class=\"text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium\">
              任务管理
            </NuxtLink>
            <NuxtLink to=\"/sprints\" class=\"text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium\">
              Sprint管理
            </NuxtLink>
            <NuxtLink to=\"/reports\" class=\"text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium\">
              报告
            </NuxtLink>
          </nav>
        </div>
      </div>
    </header>

    <!-- 主要内容 -->
    <main class=\"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8\">
      <!-- 仪表板概览 -->
      <div class=\"mb-8\">
        <h2 class=\"text-lg font-medium text-gray-900 mb-4\">项目概览</h2>
        <div class=\"grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4\">
          <!-- 统计卡片 -->
          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-gray-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M9 12h3.75M9 15h3.75M9 18h3.75m3-6.75V3.75c0-.825-.675-1.5-1.5-1.5h-6.75c0-.825-.675-1.5-1.5-1.5h-2.25c-.825 0-1.5.675-1.5 1.5v1.5c0 .825.675 1.5 1.5 1.5h2.25c.825 0 1.5-.675 1.5-1.5V3.75zm0 6.75v-3h3.75v3h-3.75z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">待办任务</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">12</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-blue-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">进行中</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">8</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-green-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">已完成</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">24</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-purple-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">当前Sprint</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">Sprint 2</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 快速操作 -->
      <div class=\"mb-8\">
        <h2 class=\"text-lg font-medium text-gray-900 mb-4\">快速操作</h2>
        <div class=\"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3\">
          <button class=\"relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2\">
            <svg class=\"mx-auto h-12 w-12 text-gray-400\" stroke=\"currentColor\" fill=\"none\" viewBox=\"0 0 48 48\">
              <path d=\"M8 14v20c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252M8 14c0 4.418 7.163 8 16 8s16-3.582 16-8M8 14c0-4.418 7.163-8 16-8s16 3.582 16 8m0 0v14m-16-4c0 4.418 7.163 8 16 8 1.38 0 2.72-.087 4-.252\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" />
            </svg>
            <span class=\"mt-2 block text-sm font-medium text-gray-900\">创建新任务</span>
          </button>

          <button class=\"relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2\">
            <svg class=\"mx-auto h-12 w-12 text-gray-400\" stroke=\"currentColor\" fill=\"none\" viewBox=\"0 0 48 48\">
              <path d=\"M8 14v20c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252M8 14c0 4.418 7.163 8 16 8s16-3.582 16-8M8 14c0-4.418 7.163-8 16-8s16 3.582 16 8m0 0v14m-16-4c0 4.418 7.163 8 16 8 1.38 0 2.72-.087 4-.252\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" />
            </svg>
            <span class=\"mt-2 block text-sm font-medium text-gray-900\">创建Sprint</span>
          </button>

          <button class=\"relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2\">
            <svg class=\"mx-auto h-12 w-12 text-gray-400\" stroke=\"currentColor\" fill=\"none\" viewBox=\"0 0 48 48\">
              <path d=\"M8 14v20c0 4.418 7.163 8 16 8 1.381 0 2.721-.087 4-.252M8 14c0 4.418 7.163 8 16 8s16-3.582 16-8M8 14c0-4.418 7.163-8 16-8s16 3.582 16 8m0 0v14m-16-4c0 4.418 7.163 8 16 8 1.38 0 2.72-.087 4-.252\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" />
            </svg>
            <span class=\"mt-2 block text-sm font-medium text-gray-900\">查看报告</span>
          </button>
        </div>
      </div>

      <!-- 最近活动 -->
      <div>
        <h2 class=\"text-lg font-medium text-gray-900 mb-4\">最近活动</h2>
        <div class=\"bg-white shadow overflow-hidden sm:rounded-md\">
          <ul role=\"list\" class=\"divide-y divide-gray-200\">
            <li class=\"px-6 py-4\">
              <div class=\"flex items-center justify-between\">
                <div class=\"flex items-center\">
                  <div class=\"flex-shrink-0 h-8 w-8\">
                    <div class=\"h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center\">
                      <span class=\"text-xs font-medium text-white\">JD</span>
                    </div>
                  </div>
                  <div class=\"ml-4\">
                    <div class=\"text-sm font-medium text-gray-900\">John Doe 创建了任务</div>
                    <div class=\"text-sm text-gray-500\">\"实现用户认证功能\" - 2小时前</div>
                  </div>
                </div>
                <div class=\"flex-shrink-0\">
                  <span class=\"inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800\">
                    新建
                  </span>
                </div>
              </div>
            </li>
            <li class=\"px-6 py-4\">
              <div class=\"flex items-center justify-between\">
                <div class=\"flex items-center\">
                  <div class=\"flex-shrink-0 h-8 w-8\">
                    <div class=\"h-8 w-8 rounded-full bg-purple-500 flex items-center justify-center\">
                      <span class=\"text-xs font-medium text-white\">AS</span>
                    </div>
                  </div>
                  <div class=\"ml-4\">
                    <div class=\"text-sm font-medium text-gray-900\">Alice Smith 完成了任务</div>
                    <div class=\"text-sm text-gray-500\">\"设计用户界面\" - 4小时前</div>
                  </div>
                </div>
                <div class=\"flex-shrink-0\">
                  <span class=\"inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800\">
                    完成
                  </span>
                </div>
              </div>
            </li>
            <li class=\"px-6 py-4\">
              <div class=\"flex items-center justify-between\">
                <div class=\"flex items-center\">
                  <div class=\"flex-shrink-0 h-8 w-8\">
                    <div class=\"h-8 w-8 rounded-full bg-orange-500 flex items-center justify-center\">
                      <span class=\"text-xs font-medium text-white\">MB</span>
                    </div>
                  </div>
                  <div class=\"ml-4\">
                    <div class=\"text-sm font-medium text-gray-900\">Mike Brown 启动了Sprint</div>
                    <div class=\"text-sm text-gray-500\">\"Sprint 2 - 高级功能\" - 1天前</div>
                  </div>
                </div>
                <div class=\"flex-shrink-0\">
                  <span class=\"inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800\">
                    Sprint
                  </span>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
// 页面元数据
useSeoMeta({
  title: 'PowerX Scrum - 仪表板',
  description: 'Scrum任务管理仪表板'
})

// 获取运行时配置
const config = useRuntimeConfig()

// 这里可以添加数据获取逻辑
// const { data: stats } = await $fetch(`${config.public.apiBaseUrl}/stats`)
</script>

<style scoped>
/* 如果需要自定义样式可以在这里添加 */
</style>