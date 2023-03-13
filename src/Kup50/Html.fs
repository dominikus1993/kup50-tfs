namespace Kup50

open HtmlDiff


module Html =
    let diff (oldHtml) (newHtml) =
        try 
            let diffHelper = HtmlDiff(oldHtml, newHtml)
            diffHelper.Build()
        with
        | _ -> newHtml

