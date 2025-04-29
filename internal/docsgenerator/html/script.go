package html

const script = `<script>
    function toggleVisibility(folderId) {
        const folder = document.getElementById(folderId);
        if (folder.classList.contains('hidden')) {
            folder.classList.remove('hidden');
        } else {
            folder.classList.add('hidden');
        }
        openFile(folderId)
    }

    function openFile(fileName) {
        fetch("/open-file?name=" + encodeURIComponent(fileName))
            .then(response => {
                if (!response.ok) {
					const errorText = response.text();
                    throw new Error(errorText || "Файл не найден или ошибка сервера");
                }
                return response.text();
            })
            .then(text => {
                document.getElementById('content').innerHTML = "<pre>" + text + "</pre>";
            })
            .catch(error => {
                document.getElementById('content').innerHTML = "<p>Ошибка: " + error.message + "</p>";
            });
    }
</script>`
