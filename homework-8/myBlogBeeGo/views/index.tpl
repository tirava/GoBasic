<!doctype html>
<html lang="ru">
<head>
    {{template "head.tpl"}}
    <title>{{.BlogName}}</title>
</head>
<body class="uk-background-muted">
<div class="uk-height-small uk-flex uk-flex-center uk-flex-middle uk-background-cover uk-light"
     data-src="static/img/backheader.jpg" uk-img>
    <h1 class="uk-heading-divider uk-heading-line uk-text-center uk-text-primary">{{.BlogName}}</h1>
</div>
{{$uname:=.UserName}}
<div class="uk-flex uk-flex-right">
    {{if $uname}}
        <div class="uk-margin-top uk-margin-right uk-flex uk-flex-middle">
            Вы вошли как:&nbsp<strong>{{$uname}}</strong>
        </div>
    {{end}}
    <button class="uk-button uk-button-default uk-margin-right uk-margin-top" type="button"
            {{if $uname}}
            onclick="document.cookie = 'myBlogBeeGo=; expires=Thu, 01 Jan 1970 00:00:00 UTC'; history.go(0)">
        Выйти
        {{else}}
            uk-toggle="target: #form-login">
            Войти
        {{end}}
    </button>
    <div id="form-login" uk-modal>
        <div class="uk-width-1-1">
            <div class="uk-container">
                <div class="uk-grid-margin uk-grid uk-grid-stack" uk-grid>
                    <div class="uk-width-1-1@m">
                        <div class="uk-margin uk-width-large uk-margin-auto uk-card uk-card-default uk-card-body uk-box-shadow-large">
                            <button class="uk-modal-close-outside" type="button" uk-close></button>
                            <h3 class="uk-card-title uk-text-center">Добро пожаловать!</h3>
                            <form>
                                <div class="uk-margin">
                                    <div class="uk-inline uk-width-1-1">
                                        <span class="uk-form-icon" uk-icon="icon: user"></span>
                                        <input id="username" class="uk-input uk-form-large" type="text"
                                               placeholder="user">
                                    </div>
                                </div>
                                <div class="uk-margin">
                                    <div class="uk-inline uk-width-1-1">
                                        <span class="uk-form-icon" uk-icon="icon: lock"></span>
                                        <input id="password" class="uk-input uk-form-large" type="password"
                                               placeholder="123">
                                    </div>
                                </div>
                                <div class="uk-margin">
                                    <button onclick="getUser();"
                                            class="uk-button uk-button-primary uk-button-large uk-width-1-1 uk-modal-close">
                                        Войти
                                    </button>
                                </div>
                                <div class="uk-text-small uk-text-center">
                                    Не зарегистрированы? <a onclick="createUser()" href="#">Создать эккаунт</a>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {{if $uname}}
        <a uk-tooltip="title: Добавить новый; pos: left" class="uk-icon-button uk-margin-top uk-margin-right"
           href="/posts/create" uk-icon="icon: plus" style="background: forestgreen; color: white;"></a>
    {{end}}
</div>

<div class="uk-margin-top uk-margin-right uk-margin-left uk-child-width-1-1 uk-grid-collapse"
     uk-height-match=".uk-card-body" uk-grid>
    {{range .Posts}}
        <div class="uk-card uk-card-default uk-card-hover uk-margin-bottom ">
            <div class="uk-card-header">
                <div class="uk-flex uk-flex-right">
                    {{if $uname}}
                        {{template "tools.tpl" .}}
                    {{end}}
                </div>
                <div class="uk-grid-small uk-flex-middle" uk-grid>
                    <div class="uk-width-auto">
                        <img class="uk-border-circle" width="40" height="40" src="../static/img/avatar.png">
                    </div>
                    <div class="uk-width-expand">
                        <h2 class="uk-card-title ">{{.Title}}</h2>
                        <p class="uk-text-meta">{{.Date2Norm}}</p>
                    </div>
                </div>
            </div>
            <div class="uk-card-body">
                <p>{{.Summary}}</p>
            </div>
            <div class="uk-card-footer uk-flex uk-flex-between">
                <a class="uk-link-heading uk-text-primary uk-button-text" href="/posts/?id={{.ID}}">Читать полностью
                    >></a>
            </div>
        </div>
    {{end}}
</div>
{{template "footer.tpl"}}
</body>
</html>


