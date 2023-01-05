/**
 * @author xiangqian
 * @date 22:24 2022/12/25
 */
;
Custom = function () {
    let obj = {}

    obj.reload = function (text) {
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

    obj.ajaxFormData = function (url, method, data, async, success, error) {
        // console.log(method, url, formData)
        // application/x-www-form-urlencoded
        $.ajax({
            url: url,
            type: method,
            data: data,
            processData: false,
            contentType: false,
            async: async,
            success: function (resp) {
                if (success) {
                    success(resp)
                }
            },
            error: function (e) {
                if (error) {
                    error(e)
                }
            }
        })
    }

    obj.ajaxJsonData = function (url, method, data, async, success, error) {
        $.ajax({
            url: url,
            type: method,
            data: data,
            dataType: "json",
            async: async,
            success: function (resp) {
                if (success) {
                    success(resp)
                }
            },
            error: function (e) {
                if (error) {
                    error(e)
                }
            }
        })
    }

    return obj
}()
;

(function () {
    function ajaxFormData(url, method, data) {
        Custom.ajaxFormData(url, method, data, true, function (resp) {
            Custom.reload(resp)
        }, function (e) {
            alert(e)
        })
    }

    let request = function ($e) {
        let formData = null
        let flag = true

        // pre func
        let pre = $e[0]._pre_
        if (pre) {
            let r = pre($e)
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
            let message = $e.attr("confirm")
            if (message) {
                flag = confirm(message)
            }
        }

        // ajaxFormData
        if (flag) {
            let url = $e.attr("href")
            // console.log(url)
            if (!(url)) {
                url = $e.attr("action")
            }
            // console.log(url)
            let method = $e.attr("method").trim().toUpperCase()
            ajaxFormData(url, method, formData)
        }
    }

    // <a></a>
    let aarr = $('a[method]')
    for (let i = 0, l = aarr.length; i < l; i++) {
        let $a = $(aarr[i])
        // console.log($a)
        $a.click(function () {
            request($a)

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
                        ajaxFormData(url, method, formData)
                        return false
                    })
                }
                break
            }
        }
    }

})()
;
