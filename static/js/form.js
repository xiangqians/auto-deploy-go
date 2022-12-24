(function () {
    let inputs = $('input[type="submit"]')
    for (let i = 0, l = inputs.length; i < l; i++) {
        let $input = $(inputs[i])
        for (let $parent = $input.parent(); !$parent.is('body'); $parent = $parent.parent()) {
            if ($parent.is('form')) {
                let $form = $parent
                let method = $form.attr("method").trim().toUpperCase()
                if (method === "POST") {
                    continue
                }

                $input.click(function () {
                    let url = $form.attr("action")
                    let method = $form.attr("method").trim().toUpperCase()
                    let formData = new FormData()
                    $form.serializeArray().forEach(e => {
                        formData.append(e.name, e.value);
                    })

                    // application/x-www-form-urlencoded
                    $.ajax({
                        url: url,
                        type: method,
                        data: formData,
                        processData: false,
                        contentType: false,
                        success: function (resp) {
                            // 使用 document.write() 覆盖当前文档
                            document.write(resp)
                            document.close()

                            // 修改当前浏览器地址
                            $span = $('span[name="_url_"]')
                            if ($span.length > 0) {
                                history.replaceState(undefined, undefined, $span.attr('url'))
                            }
                        },
                        error: function (resp) {
                            console.error(resp)
                        }
                    });

                    return false
                })
                break
            }
        }
    }
})()