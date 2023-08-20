document.addEventListener("DOMContentLoaded", function () {
    // Elements
    const monitoringFrequency = document.getElementById("monitoring-frequency");
    const notifyComments = document.getElementById("notify-comments");
    const notifyTasks = document.getElementById("notify-tasks");
    const notifyStatus = document.getElementById("notify-status");
    const saveButton = document.getElementById("save-button");
    const lastUpdate = document.getElementById("last-update");
    const activePRComments = document.getElementById("active-pr-comments");
    const activePRTasks = document.getElementById("active-pr-tasks");
    const manualButton = document.getElementById("update-button");
    const testButton = document.getElementById("test-button");

    // Save Configuration
    saveButton.addEventListener("click", function () {
        const frequencyValue = monitoringFrequency.value;
        const commentsChecked = notifyComments.checked;
        const tasksChecked = notifyTasks.checked;
        const statusChecked = notifyStatus.checked;

        console.log("Frequency:", frequencyValue);
        console.log("Notify on Comments:", commentsChecked);
        console.log("Notify on Tasks:", tasksChecked);
        console.log("Notify on Status Changes:", statusChecked);
    });

    function updateStatistics() {
        updateStatsFromBackend()
            .then(data => {
                if (!!data) {
                    lastUpdate.textContent = new Date(data.lastUpdate).toLocaleString();
                    activePRComments.textContent = data.numberOfActivePrComments;
                    activePRTasks.textContent = data.numberOfActivePrTasks;
                }
            })

    }

    updateStatistics();

    setInterval(updateStatistics, 30000);

    manualButton.addEventListener("click", function () {
        updateStatistics()
    });

    testButton.addEventListener("click", function () {
        fetch("/send-notification", {
            method: "POST",
        })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
            })
            .catch(error => {
                console.error("Error:", error);
            });
    });
});

function updateStatsFromBackend() {
    return fetch("/manual-update", {
        method: "POST",
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            if (!!data) {
                return data;
            } else {
                alert("Response was empty, but without error");
                return null;
            }
        })
        .catch(error => {
            console.error("Error:", error);
            return null;
        });

}
