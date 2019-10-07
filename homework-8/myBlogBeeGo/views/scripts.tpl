<script>
    let form = document.querySelector('form');
    if (form !== null) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    function createUser() {
        commonCreateGetUser()
    }

    function getUser() {
        commonCreateGetUser("/")
        // window.location.href = window.location.href;
        window.location.replace("/posts");
    }

    function commonCreateGetUser(id) {
        let userData = document.querySelector(`div[id="form-login"]`);
        let uname = userData.querySelector('input[id="username"]').value;
        let upass = userData.querySelector('input[id="password"]').value;
        if (uname === "" || upass === "") {
            alert("Заполните все поля!")
            return
        }
        let suffix = "";
        if (id === "/") {
            suffix = id + uname
        } else {
            alert("Отлично, теперь можно заходить!")
        }
        fetch(`/api/v1/users` + suffix, {
            method: 'POST',
            body: JSON.stringify({
                uname,
                upass,
            })
        })
    }

    function deletePost(id) {
        fetch(`/api/v1/posts/${id}`, {method: 'DELETE'})
            .then(resp => {
                window.location = "/"
            })
    }

    function updatePost(id) {
        commonCreateUpdatePost(id, `/api/v1/posts/`, 'PUT')
    }

    function createPost() {
        commonCreateUpdatePost("", `/api/v1/posts`, 'POST')
    }

    function commonCreateUpdatePost(id, api, method) {
        let postEdit = document.querySelector(`fieldset[id="editPost"]`);
        let title = postEdit.querySelector('input[id="title"]').value;
        // let date = postEdit.querySelector('input[id="date"]').value;
        let summary = postEdit.querySelector('textarea[id="summary"]').value;
        let body = postEdit.querySelector('textarea[id="body"]').value;
        fetch(api + id, {
            method: method,
            body: JSON.stringify({
                title,
                // date,
                summary,
                body,
            })
        })
            .then(resp => {
                window.location = `/posts/?id=${id}`
            })
    }
</script>
