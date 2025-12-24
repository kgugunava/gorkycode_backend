document.addEventListener("DOMContentLoaded", async () => {
  const token = localStorage.getItem("auth_token");
  console.log("Токен:", token);
  const path = location.pathname;
  console.log("Текущий путь:", path);

  // === PERSONAL ACCOUNT PAGE ===
  if (path.includes("personal") || path.includes("account")) {
    console.log("Обнаружена страница личного кабинета");
    const loader = document.getElementById("loader");
    let user;

    try {
      const res = await fetch("/api/v1/profile", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error(`Ошибка ${res.status}`);
      user = await res.json();
    } catch (err) {
      console.warn("⚠️ Сервер не отвечает", err);
    } finally {
      if (loader) loader.classList.add("hidden");
    }

    document.querySelector(".user-name").textContent = user.name;
    document.querySelector(".stats").innerHTML = `
      Сгенерировано маршрутов: ${user.routes_number ?? 0}<br><br>
      Добавлено в избранное: ${user.favourite_routes ?? 0}
    `;
  }

  // === FAVORITE DEST PAGE ===
if (location.pathname.includes("/html/favourite_dest.html")) {
  console.log("🗺️ Загружаем сохранённые маршруты...");
  const loader = document.getElementById("loader");
  const routesList = document.getElementById("routes-list");
  const modal = document.getElementById("routeModal");
  const modalTitle = modal.querySelector(".modal-title");
  const modalDescription = modal.querySelector(".modal-description");
  const mapContainer = document.getElementById("modal-map");
  const closeModal = modal.querySelector(".close-modal");

  closeModal.onclick = () => modal.classList.add("hidden");

  const token = localStorage.getItem("auth_token");
  let data;

  try {
    const res = await fetch("/api/v1/route/favourites", {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) throw new Error(`Ошибка ${res.status}`);
    data = await res.json();
  } catch (err) {
    console.warn("⚠️ Сервер не отвечает, используем тестовые данные:", err);
    data = {
      favourites: 3,
      user_id: 10,
      routes: [
        {
          route_id: 28,
          user_id: 10,
          query: {
            interests: "кофейни и музеи",
            coordinates: [56.321177, 43.990107],
            time_for_route: 90
          },
          route: {
            time: 120,
            description: "Пешеходный маршрут по центру Нижнего Новгорода",
            count_places: 2,
            places: [
              {
                title: "UserPoint",
                addres: "ул. Большая Покровская, 43",
                coordinate: [56.310043, 44.001603],
                description: "Арт-пространство и кофе-бар."
              },
              {
                title: "Галерея 9Б",
                addres: "Октябрьская улица, 9Б",
                coordinate: [56.321791, 44.00199],
                description: "Галерея современного искусства."
              }
            ]
          },
          description: "Культурно-кофейный маршрут по центру",
          is_favourite: true
        },
        {
          route_id: 31,
          user_id: 10,
          query: {
            interests: "парки и набережные",
            coordinates: [56.3269, 44.0059],
            time_for_route: 60
          },
          route: {
            time: 80,
            description: "Прогулка вдоль реки и через парки",
            count_places: 3,
            places: [
              {
                title: "Александровский сад",
                addres: "ул. Минина и Пожарского, 2Б",
                coordinate: [56.3291, 44.0073],
                description: "Исторический городской парк."
              },
              {
                title: "Чкаловская лестница",
                addres: "ул. Верхне-Волжская набережная",
                coordinate: [56.3299, 44.0108],
                description: "Знаковая лестница с видом на Волгу."
              },
              {
                title: "Набережная Федоровского",
                addres: "наб. Федоровского, 15",
                coordinate: [56.3328, 44.0191],
                description: "Панорамная прогулка с видом на реку."
              }
            ]
          },
          description: "Природный маршрут по набережной",
          is_favourite: true
        },
        {
          route_id: 42,
          user_id: 10,
          query: {
            interests: "еда и развлечения",
            coordinates: [56.325, 44.02],
            time_for_route: 100
          },
          route: {
            time: 110,
            description: "Маршрут для любителей вкусно поесть",
            count_places: 3,
            places: [
              {
                title: "Traveler's Coffee",
                addres: "ул. Рождественская, 25",
                coordinate: [56.3284, 44.0081],
                description: "Кофейня с уютной атмосферой."
              },
              {
                title: "Hurma Bar",
                addres: "ул. Большая Покровская, 65",
                coordinate: [56.3132, 44.0045],
                description: "Бар с авторскими коктейлями."
              },
              {
                title: "Кинотеатр «Орленок»",
                addres: "ул. Пискунова, 11",
                coordinate: [56.3206, 44.0065],
                description: "Классический кинотеатр с ретро-залом."
              }
            ]
          },
          description: "Вечерний гастро-маршрут",
          is_favourite: true
        }
      ]
    };
    
  } finally {
    loader.classList.add("hidden");
  }

  const favRoutes = data.routes?.filter(r => r.is_favourite && r.route?.places?.length > 0) ?? [];

  // if (favRoutes.length === 0) {
  //   routesList.innerHTML = "<p>У вас пока нет сохранённых маршрутов.</p>";
  //   return;
  // }

  favRoutes.forEach((route, i) => {
    const places = route.route.places;
    const start = places[0];
    const end = places[places.length - 1];

    const startAddress = start?.addres || start?.title || "Неизвестно";
    const endAddress = end?.addres || end?.title || "Неизвестно";

    const card = document.createElement("div");
    card.className = "route-card";
    const title = `${startAddress} — ${endAddress}`;
    card.innerHTML = `
        <h3>${title}</h3>
        <button class="open-route-btn">Открыть</button>
    `;


    card.querySelector(".open-route-btn").addEventListener("click", () => showRouteModal(route));
    routesList.appendChild(card);
  });

  function showRouteModal(route) {
    modal.classList.remove("hidden");
    const start = route.route.places[0];
    const end = route.route.places[route.route.places.length - 1];
    const startTitle = start?.addres || start?.title || "Неизвестно";
    const endTitle = end?.addres || end?.title || "Неизвестно";
    modalTitle.textContent = `${startTitle} — ${endTitle}`;

    mapContainer.innerHTML = "";

    if (typeof ymaps === "undefined") {
      mapContainer.innerHTML = "<p>⚠️ Карта не загрузилась.</p>";
      return;
    }

    ymaps.ready(() => {
      const startCenter = route.query?.coordinates || route.route.places[0]?.coordinate || [56.3269, 44.0059];
      const map = new ymaps.Map(mapContainer, { center: startCenter, zoom: 12 });

      const coords = [];
      route.route.places.forEach(p => {
        if (p.coordinate && Array.isArray(p.coordinate)) {
          coords.push(p.coordinate);
          const placemark = new ymaps.Placemark(p.coordinate, {
            hintContent: p.title || "Точка маршрута",
            balloonContent: `<b>${p.title || "Без названия"}</b><br>${p.addres || ""}<br>${p.description || ""}`
          }, {
            iconLayout: 'default#image',                  
            iconImageHref: '/static/src/images/point.png',      
            iconImageSize: [40, 60],                      
            iconImageOffset: [-20, -40],                   
            draggable: true              
          });
          map.geoObjects.add(placemark);
        }
      });

      if (coords.length > 1) {
        const line = new ymaps.Polyline(coords, {}, {
          strokeColor: "#FF6600",
          strokeWidth: 4
        });
        map.geoObjects.add(line);
        map.setBounds(map.geoObjects.getBounds(), { checkZoomRange: true, zoomMargin: 30 });
      }
    });
  }
}
});
