document.addEventListener("DOMContentLoaded", () => {
    mountCopyWidget()
})

function mountCopyWidget() {
    const copyWidgetElements = Array.from(document.querySelectorAll("[data-copy-element]"))
    const copyWidgets = copyWidgetElements.map (element => ({
        trigger: element,
        target: document.querySelector(element.dataset.copyElement),
        valueToCopyFromTarget: element.dataset.copyProp
    }))

    const validCopyWidgets = copyWidgets.filter (({ target, valueToCopyFromTarget }) => !!target && !!valueToCopyFromTarget)
    const invalidCopyWidgets = copyWidgets.filter (({ target, valueToCopyFromTarget }) => !target || !valueToCopyFromTarget)

    for (const widget of invalidCopyWidgets) {
        console.log("Copy widget either miss its target or the name of the property to copy from the target", widget)
    }

    for (const { target, trigger, valueToCopyFromTarget } of validCopyWidgets) {
        trigger.addEventListener("click", () => {
            const textToCopy = target[valueToCopyFromTarget]
            navigator.clipboard.writeText(textToCopy)
            trigger.innerHTML = `&check;${trigger.innerHTML}`
            console.log("Text copied!", textToCopy)
        })
    }
}