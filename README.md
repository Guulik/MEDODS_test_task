# Сервис авторизации

## Отправка предупреждения на email
Решил использовать MailHog как SMTP сервер для локального тестирования. Сообщения приходят :)
![image](https://github.com/user-attachments/assets/2f3178cc-908c-456d-8d0e-2429bc8f3538)
В сообщении указан конкретный ip, с которого был выполнен рефреш. Ip, равно как и Email получателя моковые
![image](https://github.com/user-attachments/assets/a53b23fa-9d68-4544-9fac-2ae3845787e1)

## Передача подписи jwt
Секретный ключ jwt я записываю в Docker secret, прокидываю его через compose и корректно достаю из приложения. Всё для безопасности ключа.
А потом в мейк файле явно его создаю, чем невелирую все предпринятые шаги по обеспечению безопасонсти :))

В реальном проекте я бы так не поступал, но здесь это сделано исключительно ради удобства и ускорения проверки.

# WIP
проект на стадии полировки и тестирования. Сейчас есть все основные сервисы и функции. Осталось покрыть тестами, протестировать руками и дописать readme.
