<!doctype html>
<html lang="ru">
<head>
    {{template "head.tpl"}}
    <title>Блог Евгения Климова</title>
</head>
<body class="uk-background-muted">
<div class="uk-height-small uk-flex uk-flex-center uk-flex-middle uk-background-cover uk-light"
     data-src="static/images/backheader.jpg" uk-img>
    <h1 class="uk-heading-divider uk-heading-line uk-text-center uk-text-primary">Блог Евгения Климова</h1>
</div>

<div class="uk-flex uk-flex-right">
    <a uk-tooltip="title: Добавить новый; pos: left" class="uk-icon-button uk-margin-top uk-margin-right"
       href="/posts/create" uk-icon="icon: plus" style="background: forestgreen; color: white;"></a>
</div>

<div class="uk-margin-top uk-margin-right uk-margin-left uk-child-width-1-1 uk-grid-collapse"
     uk-height-match=".uk-card-body" uk-grid>
    {{range .}}
        <div class="uk-card uk-card-default uk-card-hover uk-margin-bottom ">
            <div class="uk-card-header">
                <div class="uk-flex uk-flex-right">
                    {{template "tools.tpl" .}}
                </div>
                <div class="uk-grid-small uk-flex-middle" uk-grid>
                    <div class="uk-width-auto">
                        <img class="uk-border-circle" width="40" height="40" src="../static/images/avatar.png">
                    </div>
                    <div class="uk-width-expand">
                        <h2 class="uk-card-title ">{{.Title}}</h2>
                        <p class="uk-text-meta">{{.Date}}</p>
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


