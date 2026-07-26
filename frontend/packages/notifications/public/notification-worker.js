self.addEventListener("install", () => {
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});

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

  event.waitUntil(
    (async () => {
      const clients = await self.clients.matchAll({ type: "window", includeUncontrolled: true })
      for (const client of clients) {
        client.postMessage({ type: "push", title, body: options.body })
      }
      const bc = new BroadcastChannel("order-updates")
      bc.postMessage({ type: "push", title, body: options.body })
      bc.close()
      return self.registration.showNotification(title, options)
    })(),
  )
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
