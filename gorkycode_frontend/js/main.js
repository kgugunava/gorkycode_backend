// === Multi-step init + управление шагами ===
const container = document.querySelector(".multi-step-container");
const circles = container.querySelectorAll(".circle");
const progressBar = container.querySelector(".indicator");
const buttons = container.querySelectorAll(".buttons button"); // явнее получаем Prev/Next
const contents = container.querySelectorAll(".content");

let currentStep = 1;
const totalSteps = circles.length;

const updateSteps = (e) => {
  const isNext = e.target.id === "next";
  const isPrev = e.target.id === "prev";

  // Безопасное изменение currentStep в границах [1, totalSteps]
  if (isNext && currentStep < totalSteps) {
    // Проверка перед переходом на второй шаг
    if (currentStep === 1) {
      const startPoint = localStorage.getItem("startPoint");
      if (!startPoint) {
        alert("Выберите начальную точку на карте!");
        return; // останавливаем переход
      }
    }
    currentStep++;
  }
  if (isPrev && currentStep > 1) currentStep--;

  // обновляем кружки
  circles.forEach((circle, index) => {
    circle.classList.toggle("active", index < currentStep);
  });

  // обновляем прогресс-бар
  progressBar.style.width = `${((currentStep - 1) / (totalSteps - 1)) * 100}%`;

  // плавное скрытие предыдущего блока
  contents.forEach((block, index) => {
    if (index + 1 === currentStep) {
      block.classList.remove("fade-out");
      block.classList.add("active");
    } else {
      // чтобы корректно убрать active и предотвратить конфликт анимаций
      block.classList.remove("active");
      block.classList.remove("fade-out");
      // если это не активный — делаем fade-out кратковременно (опц.)
      if (index + 1 < currentStep) {
        block.classList.add("fade-out");
      }
    }
  });

  // блокировка кнопок: Prev выключен на первом, Next выключен на последнем
  const prevBtn = buttons[0];
  const nextBtn = buttons[1];
  prevBtn.disabled = currentStep === 1;
  nextBtn.disabled = currentStep === totalSteps;

  // если мы на последнем шаге — дополнительно снимаем/добавляем кастомные классы при необходимости
  // (например, можно подсветить финальную секцию)
};

// Повесим слушатели на Prev/Next (как раньше)
buttons.forEach((btn) => btn.addEventListener("click", updateSteps));

// === Обработка скролла между секциями (с уже имеющейся логикой) ===
let isScrolling = false;

function scrollToSection(sectionId) {
  const element = document.getElementById(sectionId);
  if (element) {
    isScrolling = true;
    element.scrollIntoView({ behavior: "smooth" });

    // Разрешаем скролл через 900ms (как было)
    setTimeout(() => {
      isScrolling = false;
    }, 900);
  }
}

// Блокируем обычный скролл колесиком
window.addEventListener('wheel', (e) => {
  if (isScrolling) {
    e.preventDefault();
    return;
  }

  const sections = document.querySelectorAll('.section');
  const currentScroll = window.scrollY;
  const windowHeight = window.innerHeight;

  let currentSectionIndex = -1;
  sections.forEach((section, index) => {
    const rect = section.getBoundingClientRect();
    if (rect.top >= 0 && rect.top < windowHeight / 2) {
      currentSectionIndex = index;
    }
  });

  if (currentSectionIndex !== -1) {
    if (e.deltaY > 0 && currentSectionIndex < sections.length - 1) {
      e.preventDefault();
      scrollToSection(sections[currentSectionIndex + 1].id);
    } else if (e.deltaY < 0 && currentSectionIndex > 0) {
      e.preventDefault();
      scrollToSection(sections[currentSectionIndex - 1].id);
    }
  }
}, { passive: false });

// Блокируем скролл тачпадом
window.addEventListener('touchmove', (e) => {
  if (isScrolling) {
    e.preventDefault();
  }
}, { passive: false });

// === Кнопка "Сгенерируй свой путь" (переход вниз) ===
const generateBtn = document.querySelector(".generate-btn");
if (generateBtn) {
  generateBtn.addEventListener("click", (ev) => {
    ev.preventDefault();
    scrollToSection("section2");
  });
}

// === Логика для стрелки вниз/вверх ===
const scrollArrowWrapper = document.querySelector(".scroll-arrow");
const scrollArrowImg = scrollArrowWrapper ? scrollArrowWrapper.querySelector("img") : null;
let arrowDown = true;

if (scrollArrowWrapper) {
  scrollArrowWrapper.addEventListener("click", (ev) => {
    ev.preventDefault();
    const sections = document.querySelectorAll(".section");
    let currentSectionIndex = -1;

    sections.forEach((section, index) => {
      const rect = section.getBoundingClientRect();
      if (rect.top >= 0 && rect.top < window.innerHeight / 2) {
        currentSectionIndex = index;
      }
    });

    if (arrowDown && currentSectionIndex < sections.length - 1) {
      scrollToSection(sections[currentSectionIndex + 1].id);
    } else if (!arrowDown && currentSectionIndex > 0) {
      scrollToSection(sections[currentSectionIndex - 1].id);
    }

    arrowDown = !arrowDown;
    if (scrollArrowImg) {
      scrollArrowImg.style.transform = arrowDown ? "rotate(0deg)" : "rotate(180deg)";
      scrollArrowImg.style.transition = "transform 0.4s ease";
    }
  });
}

const routeData = {
    coordinates: null,     // [lat, lon]
    time_for_route: null,  // минуты (число)
    interests: null        // текст
};

// Запускаем инициализацию карты при загрузке страницы
window.addEventListener('load', () => {
  initMap();
});

// ----------------------
// === Сбор опций (time, interests) в routeData в реальном времени ===
// ----------------------
const freeTimeSelect = document.getElementById('free_time');
if (freeTimeSelect) {
  freeTimeSelect.addEventListener('change', (e) => {
    const v = parseInt(e.target.value, 10);
    routeData.time_for_route = Number.isFinite(v) ? v : null;
  });
}

const userWishes = document.getElementById('user_wishes');
if (userWishes) {
  userWishes.addEventListener('input', (e) => {
    routeData.interests = e.target.value.trim() || null;
  });
}

// ----------------------
// === Сохранение в localStorage ===
// ----------------------
const timeSelect = document.getElementById("free_time");
if (timeSelect) {
  timeSelect.addEventListener("change", () => {
    const selectedTime = timeSelect.value;
    localStorage.setItem("selectedTime", selectedTime);
    routeData.time_for_route = parseInt(selectedTime, 10);
    console.log("Выбрано время:", selectedTime);
  });
}

const wishesInput = document.getElementById("user_wishes");
if (wishesInput) {
  wishesInput.addEventListener("input", () => {
    const interests = wishesInput.value.trim();
    localStorage.setItem("userInterests", interests);
    routeData.interests = interests;
    console.log("Интересы:", interests);
  });
}

const mapContainer = document.getElementById("yandex-map");
if (mapContainer) {
  mapContainer.addEventListener("mouseleave", () => {
    mapContainer.style.pointerEvents = "none";
    setTimeout(() => (mapContainer.style.pointerEvents = "auto"), 500);
  });
}

// ----------------------
// === Отправка JSON на сервер при нажатии "Создать маршрут" ===
// ----------------------
// const createRouteBtn = document.getElementById("createRouteBtn");
// if (createRouteBtn) {
//   createRouteBtn.addEventListener("click", async () => {
//     // Собираем все данные из localStorage
//     const startPoint = localStorage.getItem("startPoint");
//     const freeTime = localStorage.getItem("selectedTime");
//     const userInterests = localStorage.getItem("userInterests");

//     // Проверим, что данные есть
//     if (!startPoint) {
//       alert("Выберите начальную точку на карте!");
//       return;
//     }
//     if (!freeTime) {
//       alert("Укажите, сколько у вас свободного времени!");
//       return;
//     }

//     // Формируем JSON с новой структурой
//     const routeData = {
//       coordinates: JSON.parse(startPoint),
//       time_for_route: parseInt(freeTime, 10),
//       interests: userInterests || ""
//     };

//     console.log("📦 Отправляем маршрут:", routeData);

//     // Отправка на сервер
//     try {
//       const response = await fetch("http://localhost:8020/api/v1/create-route", {
//         method: "POST",
//         headers: { "Content-Type": "application/json",
//                     "Authorization": "Bearer " + localStorage.getItem("auth_token")
//                   },
//         body: JSON.stringify(routeData),
//       });

//       if (!response.ok) throw new Error(`Ошибка: ${response.status}`);

//       const result = await response.json();
//       console.log("✅ Сервер вернул:", result);

//       alert("Маршрут успешно создан!");
//     } catch (err) {
//       console.error("❌ Ошибка при отправке:", err);
//       alert("Произошла ошибка при создании маршрута.");
//     }
//   });
// }

// === Отображение выбранных данных в финальном протоколе ===
function updateFinalProtocol() {
  const finalPoint = document.getElementById("final-point");
  const finalTime = document.getElementById("final-time");
  const finalWishes = document.getElementById("final-wishes");

  const startPoint = localStorage.getItem("startPoint");
  const selectedTime = localStorage.getItem("selectedTime");
  const userInterests = localStorage.getItem("userInterests");

  // Координаты точки
  if (startPoint) {
    const coords = JSON.parse(startPoint);
    finalPoint.textContent = `[${coords[0].toFixed(5)}, ${coords[1].toFixed(5)}]`;
  } else {
    finalPoint.textContent = "не выбрано";
  }

  // Время
  if (selectedTime) {
    const minutes = parseInt(selectedTime, 10);
    let timeStr = "";
    if (minutes < 60) timeStr = `${minutes} минут`;
    else if (minutes === 60) timeStr = "1 час";
    else if (minutes === 90) timeStr = "1.5 часа";
    else if (minutes === 120) timeStr = "2 часа";
    else timeStr = "3 часа и более";
    finalTime.textContent = timeStr;
  } else {
    finalTime.textContent = "не выбрано";
  }

  // Интересы
  if (userInterests && userInterests.trim() !== "") {
    finalWishes.textContent = userInterests;
  } else {
    finalWishes.textContent = "не выбрано";
  }
}

// Когда пользователь переходит на последний шаг (4-й)
document.getElementById("next").addEventListener("click", () => {
  const circles = document.querySelectorAll(".circle");
  const activeCount = document.querySelectorAll(".circle.active").length;
  if (activeCount === circles.length) {
    updateFinalProtocol();
  }
});

window.addEventListener("load", () => {
  localStorage.removeItem("startPoint");
  localStorage.removeItem("selectedTime");
  localStorage.removeItem("userInterests");
  console.log("🔄 Очищены старые параметры маршрута");
});

// === Эффект печатания текста ===
document.addEventListener("DOMContentLoaded", () => {
  const el = document.getElementById("typewriter");
  const text = el.textContent;
  el.textContent = ""; // очищаем, чтобы печатать заново

  let i = 0;
  const speed = 50; // скорость (мс на символ)

  function type() {
    if (i < text.length) {
      el.textContent += text.charAt(i);
      i++;
      setTimeout(type, speed);
    }
  }

  type();
});