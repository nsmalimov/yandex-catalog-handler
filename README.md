TODO:

Есть файлы Прайсов в формате YML (Яндекс маркет). Каждый файл весит по 300мб. В каждом примерно по 300тыс. товаров. Всего таких прайса четыре. Возможно увеличение до 6.

В этих файлах есть дублирующиеся позиции товарных предложений (совпадает пара артикул+бренд). Дубли необходимо удалить, оставив только уникальные предложения суммарно по всем четырём файлам согласно паре артикул + бренд. Должно быть сравнение между собой.

Исходники файлов доступны по прямым ссылкам без авторизации. После того как прайсы будут очищены от дублей, необходимо выгружать их на хостинг, чтобы они были доступны по прямым ссылкам аналогичным образом. Ссылки в дальнейшем прикрепляю к агрегатору, для обновления цен у товарных предложений. Хостинг предоставлю, нужно лишь сказать какие должны быть у него технические требования.

Обновление прайсов происходит 1 раз в день, как следствие все это должно работать автономно каждый день. Админка не нужна.  Настройки могу делать в исходниках, но логи хотелось бы видеть с историей и результатами обновлений.

Есть ещё параметр «цена» в прайсе, в идеале бы оставлять наибольшую цену. То есть перед тем как удалить дубли, необходимо сверить условие и удалить в данном случае тогда всё кроме предложения с наибольшей ценой в паре артикул+бренд

Так же из пожеланий: в логах надо видеть сколько было изначально строк со всеми товарными предложения, сколько дублей убрано и конечное количество уникальных предложений. Все это по каждому файлу в итоге.

// develop

79.143.31.238

data/30616439-6263-3866-2D35-6364632D3131&FranchiseeId=383450

http://5.101.51.209/data/66343037-3430-3935-2D35-3163632D3131&FranchiseeId=383450
http://5.101.51.209/data/30616439-6263-3866-2D35-6364632D3131&FranchiseeId=383450
http://5.101.51.209/data/36366335-6664-3661-2D64-3761342D3131&FranchiseeId=383450
http://5.101.51.209/data/65303631-3762-6565-2D64-6335362D3131&FranchiseeId=383450

docker run -d -p 5432:5432 --name my-postgres -e POSTGRES_PASSWORD=123 postgres

docker run -p 8080:8896 -d --name yandex-catalog-handler yandex-catalog-handler

docker run -p 8080:8896 -d -it --name yandex-catalog-handler yandex-catalog-handler

docker build -t yandex-catalog-handler .

docker run -p 8080:8896 -it --memory="4g" --memory-swap="3g" -v /Users/nurislam_alimov/IdeaProjects/yandex-catalog-handler/data:/app/data --name yandex-catalog-handler yandex-catalog-handler

docker run -p 8080:8896 -d -v /var/www/app/static/data:/app/data --name yandex-catalog-handler yandex-catalog-handler

go tool pprof -gif http://localhost:6060/debug/pprof/profile

Need

логи с пагинацией

запуск с кнопки

расписание запуска

добавление еще 1 файла (не прожует)

анализ и понимание где утечка памяти