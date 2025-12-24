 // === Инициализация Яндекс.Карты ===
 let map, placemark;

 function initMap() {
   map = new ymaps.Map("yandex-map", {
     center: [56.3269, 44.0059], // Нижний Новгород
     zoom: 12,
   });
 
   // Обработка клика по карте
   map.events.add("click", function (e) {
     const coords = e.get("coords");
 
     // Удаляем предыдущую метку
     if (placemark) {
       map.geoObjects.remove(placemark);
     }
 
     // Добавляем новую метку
     //placemark = new ymaps.Placemark(coords, {}, { draggable: true });
     placemark = new ymaps.Placemark(
       coords, 
       {}, 
       {
         // ↓↓↓ ЭТО НОВЫЕ НАСТРОЙКИ ↓↓↓
         iconLayout: 'default#image',                    // Используем кастомное изображение
         iconImageHref: '/static/src/images/point.png',       // Путь к вашей PNG метке
         iconImageSize: [40, 50],                       // Размер иконки
         iconImageOffset: [-20, -40],                   // Смещение для центрирования
         draggable: true                                // Эта опция была и раньше
         // ↑↑↑ ЭТО НОВЫЕ НАСТРОЙКИ ↑↑↑
       }
     );
     map.geoObjects.add(placemark);
 
     // Сохраняем координаты
     localStorage.setItem("startPoint", JSON.stringify(coords));
 
     console.log("Выбрана точка:", coords);
   });
 }
 
 if (typeof ymaps !== "undefined") {
   ymaps.ready(initMap);
 }