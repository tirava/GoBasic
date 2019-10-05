<script>
    let form = document.querySelector('form');
    if (form !== null) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
        });
    }

    function createUser() {
        let userData = document.querySelector(`div[id="form-login"]`);
        let uname = userData.querySelector('input[id="username"]').value;
        let upass = userData.querySelector('input[id="password"]').value;
        fetch(`/api/v1/users`, {
            method: 'POST',
            body: JSON.stringify({
                uname,
                upass,
            })
        })
        alert("Отлично, теперь можно заходить!")
            // .then(resp => {
            //     window.location = "/"
            // })
    }

    function getUser() {
        let userData = document.querySelector(`div[id="form-login"]`);
        let uname = userData.querySelector('input[id="username"]').value;
        let upass = userData.querySelector('input[id="password"]').value;
        fetch(`/api/v1/users/`+uname, {
            method: 'POST',
            body: JSON.stringify({
                uname,
                upass,
            })
        })
        .then(resp => {
            window.location = "/"
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
