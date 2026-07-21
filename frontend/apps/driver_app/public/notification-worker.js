self.addEventListener("push", (event) => {
  const data = event.data?.json()

  const title = data?.title ?? "Notification"
  const options = {
    body: data?.body ?? "",
    icon: data?.icon ?? "/favicon.ico",
    badge: "/favicon.ico",
    data: {
      url: data?.url ?? "/",
    },
  }

  event.waitUntil(self.registration.showNotification(title, options))
})

self.addEventListener("notificationclick", (event) => {
  event.notification.close()

  const url = event.notification.data?.url ?? "/"

  event.waitUntil(
    self.clients.matchAll({ type: "window" }).then((clients) => {
      for (const client of clients) {
        if (client.url.includes(url) && "focus" in client) {
          return client.focus()
        }
      }
      if (self.clients.openWindow) {
        return self.clients.openWindow(url)
      }
    }),
  )
})
