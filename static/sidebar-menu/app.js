const toggleButton = document.getElementById('toggle-btn')
const sidebar = document.getElementById('sidebar')

function toggleSidebar(){
  sidebar.classList.toggle('close')
  toggleButton.classList.toggle('rotate')

  closeAllSubMenus()
}

function toggleSubMenu(button){

  if(!button.nextElementSibling.classList.contains('show')){
    closeAllSubMenus()
  }

  button.nextElementSibling.classList.toggle('show')
  button.classList.toggle('rotate')

  if(sidebar.classList.contains('close')){
    sidebar.classList.toggle('close')
    toggleButton.classList.toggle('rotate')
  }
}

function closeAllSubMenus(){
  Array.from(sidebar.getElementsByClassName('show')).forEach(ul => {
    ul.classList.remove('show')
    ul.previousElementSibling.classList.remove('rotate')
  })
}

function send(){
  const Http = new XMLHttpRequest();
  const url='https://jsonplaceholder.typicode.com/users';
  Http.open("GET", url);
  Http.send();

  Http.onreadystatechange = (e) => {

    if (this.status == 'asdf') {
      console.log(Http.responseText);
    } else {
      launch_toast()
    }

    // console.log(Http.responseText);  // => получим массив данных в формате JSON
  }
}

function send_post(){
  const user = {
    "name": "Ivan Ivanov",
    "username": "ivan2002",
    "email": "ivan2002@mail.com",
  };
  
  const Http = new XMLHttpRequest();
  const url='https://jsonplaceholder.typicode.com/users';
  Http.open("POST", url);
  Http.setRequestHeader("Content-Type", "application/json");
  Http.send(JSON.stringify(user));

  Http.upload.onprogress = function(e) {
    console.log(e.loaded);  // загруженные данные
    console.log(e.total);  // всего данных
  }
  
  // окончание загрузки данных
  Http.upload.onload = function(e) {
    console.log("Данные загружены");
  }
  
  // полная отправка запроса
  Http.onload = function() {
    console.log(Http.status);
  }
  
  // ошибка при отправке запроса
  Http.onerror = function() {
    console.log("Ошибка запроса");
  }
}