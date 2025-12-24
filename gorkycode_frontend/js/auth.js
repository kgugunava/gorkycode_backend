const loginForm = document.getElementById('loginForm');
if (loginForm) {
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('login_email').value.trim();
    const password = document.getElementById('login_password').value;

    if (!email || !password) {
      alert('Пожалуйста, заполните email и пароль.');
      return;
    }

    const payload = { email, password };

    try {
      const res = await fetch('/api/v1/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const errText = await res.text();
        alert('Ошибка входа: ' + errText);
        return;
      }

      const data = await res.json();
      // Ожидаем { token: "...", user: { ... } }
      if (data.token) {
        localStorage.setItem('auth_token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user || {}));
      }

      alert('Вход выполнен успешно');
      // При желании: вернуть на главную страницу
      location.href = '/html/index.html';
    } catch (err) {
      console.error(err);
      alert('Ошибка сети при входе.');
    }
  });
}

// registration form
const regForm = document.getElementById('registerForm');
if (regForm) {
  regForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const fullname = document.getElementById('reg_fullname').value.trim();
    const username = document.getElementById('reg_username').value.trim();
    const email = document.getElementById('reg_email').value.trim();
    const password = document.getElementById('reg_password').value;
    const password2 = document.getElementById('reg_password2').value;

    if (!email || !password || !fullname) {
      alert('Пожалуйста, заполните обязательные поля.');
      return;
    }
    if (password !== password2) {
      alert('Пароли не совпадают.');
      return;
    }

    const payload = { fullname, username, email, password };

    try {
      const res = await fetch('/api/v1/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const err = await res.text();
        alert('Ошибка регистрации: ' + err);
        return;
      }

      const data = await res.json();
      alert('Регистрация прошла успешно. Войдите в аккаунт.');
      location.href = '/html/login.html';
    } catch (err) {
      console.error(err);
      alert('Ошибка сети при регистрации.');
    }
  });
}
