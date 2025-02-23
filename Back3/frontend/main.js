const errors = [
    {
        name: "ФИО",
        description: "Поле должно содержать полное ФИО (фамилия, имя, отчество) на кириллице, " +
            "состоящее минимум из трех слов, разделенных пробелами; допустимы только буквы(например, Иванов Иван Иванович).",
    },
    {
        name: "Номер телефона",
        description: "Поле должно содержать номер телефона в формате +XXXXXXXXXXX; допустимы только цифры.",
    },
    {
        name: "Адрес электронной почты",
        description: "Поле должно содержать корректный адрес электронной почты в формате example@domain.com; допустимы буквы, цифры, точки, дефисы и символ @.",
    },
    {
        name: "Дата рождения",
        description: "Поле должно содержать дату в формате ГГГГ-ММ-ДД;",
    },
    {
        name: "Пол",
        description: "Поле обязательно для выбора; выберите один из вариантов: \"М\" (мужской) или \"Ж\" (женский).",
    },
    {
        name: "Любимый язык программирования",
        description: "Поле должно содержать хотя бы один выбранный язык программирования; можно выбрать несколько вариантов из списка.",
    },
    {
        name: "Биография",
        description: "Sample text",
    },
    {
        name: "С контрактом ознакомлен(а)",
        description: "Поле обязательно для подтверждения; поставьте галочку, чтобы подтвердить ознакомление с контрактом.",
    },
]
const backUrl = 'http://158.160.154.2:8080/process'
document.addEventListener('DOMContentLoaded', function() {
const forma = document.getElementById('forma')
forma.addEventListener('submit', function (event) {

    event.preventDefault() // отмено перезагрузко
    
    const formData = new FormData(this)
 

    const errContainer = document.getElementById('errContainer')
    const loader = document.getElementById("form-loader")
    errContainer.classList.remove('error')
    loader.style.visibility = "visible"
    loader.style.display = "block"
    
    fetch(backUrl, {
        method: 'POST',
        body: formData,
    })
        .then((response) => {
            let errorsHtml = `<ol>`
            console.log(response.status)
            switch (response.status) {
            case 200:{
                loader.style.visibility = "hidden"
                loader.style.display = "none"
                errContainer.classList.add('success')
                errContainer.innerText = 'Форма успешно отправлена!'
                break;  
            }
            case 400: {
                return response.json().then(data => {
                    console.log(data)
                    data.forEach((code, key) => {
                     errorsHtml += `<li><h6>Поле "${errors[code - 1].name}" введено не верно.</h6><p>${errors[code - 1].description}</p></li>`
                    })
                    errorsHtml += `</ol>`

                    loader.style.visibility = "hidden"
                    loader.style.display = "none"

                    errContainer.classList.remove("success")
                    errContainer.classList.add("error")
                    errContainer.innerHTML = errorsHtml
                    })
                    
                break;
            }
            default: {
                loader.style.visibility = "hidden"
                loader.style.display = "none"

                errContainer.classList.remove("success")
                errContainer.classList.add("error")
                errContainer.innerHTML = `<p>{response.status} Непредвиденная ошибка</p>`
                break;
            }
            }
        })
        .catch((error) => {
            console.error(error)

        });

});
})