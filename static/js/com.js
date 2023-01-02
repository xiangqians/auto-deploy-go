/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
(function () {

    function reload(text) {
        // 使用 document.write() 覆盖当前文档
        document.write(text)
        document.close()

        // 修改当前浏览器地址
        let $html = $('html')
        let url = $html.attr('url')
        if (url) {
            history.replaceState(undefined, undefined, url)
        }
    }

    function ajax(url, method, formData) {
        // console.log(method, url, formData)
        // application/x-www-form-urlencoded
        $.ajax({
            url: url,
            type: method,
            data: formData,
            processData: false,
            contentType: false,
            success: function (resp) {
                reload(resp)
            },
            error: function (e) {
                alert(e)
            }
        })
    }

    // <a></a>
    let aarr = $('a[method]')
    for (let i = 0, l = aarr.length; i < l; i++) {
        let $a = $(aarr[i])
        // console.log($a)
        $a.click(function () {
            let formData = null
            let flag = true

            // pre func
            let pre = $a[0]._pre_
            if (pre) {
                let r = pre()
                let rarr = null
                let rl = 0
                if (Object.prototype.toString.call(r) === '[object Boolean]') {
                    flag = r
                } else if (Object.prototype.toString.call(r) === '[object Array]' && (rl = (rarr = r).length) > 0) {
                    flag = rarr[0] ? true : false
                    if (flag && rl > 1) {
                        formData = new FormData()
                        for (let ri = 1; ri < rl; ri++) {
                            let robj = rarr[ri]
                            for (let name in robj) {
                                formData.append(name, robj[name])
                            }
                        }
                    }
                }
            }
            // confirm
            else {
                let message = $a.attr("confirm")
                if (message) {
                    flag = confirm(message)
                }
            }

            // ajax
            if (flag) {
                let url = $a.attr("href")
                let method = $a.attr("method").trim().toUpperCase()
                ajax(url, method, formData)
            }

            // 取消 <a></a> 默认行为
            return false
        })
    }

    // <form></form>
    let inputs = $('input[type="submit"]')
    for (let i = 0, l = inputs.length; i < l; i++) {
        let $input = $(inputs[i])
        for (let $parent = $input.parent(); !$parent.is('body'); $parent = $parent.parent()) {
            if ($parent.is('form')) {
                let $form = $parent
                let method = $form.attr("method").trim().toUpperCase()
                if (method !== "POST") {
                    $input.click(function () {
                        let url = $form.attr("action")
                        let method = $form.attr("method").trim().toUpperCase()
                        let formData = new FormData()
                        $form.serializeArray().forEach(e => {
                            formData.append(e.name, e.value);
                        })
                        ajax(url, method, formData)
                        return false
                    })
                }
                break
            }
        }
    }

    // header.html -> nav
    let $html = $('html')
    // console.log('$html', $html)
    let url = $html.attr('url')
    // console.log('url', url)
    let $navs = $('a[name="nav"]')
    // console.log('$navs', $navs)
    for (let i = 0, length = $navs.length; i < length; i++) {
        let $nav = $($navs[i])
        let navUrl = $nav.attr("href")
        if (url === navUrl) {
            // console.log('navUrl', navUrl)
            $nav.css('background-color', 'lightgray');
        }
    }
})()