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

  function manualUpdate() {
    getManualUpdateFromBackEnd().then((data) => {
      if (!!data) {
        lastUpdate.textContent = new Date(
          data.LastUpdate * 1000
        ).toLocaleString();
        activePRComments.textContent = data.NumberOfActivePrComments;
        activePRTasks.textContent = data.NumberOfActivePrTasks;
      }
    });
  }

  updateStatistics();

  setInterval(updateStatistics, 30000);

  manualButton.addEventListener("click", function () {
    manualUpdate();
  });

  testButton.addEventListener("click", function () {
    sendNotification();
  });
});

function sendNotification() {
  fetch("/send-notification", {
    method: "POST",
  })
    .then((response) => response.json())
    .then((data) => {
      alert(data.message);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}

async function getManualUpdateFromBackEnd() {
  try {
    const response = await fetch("/manual-update", {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error("Network response was not ok: " + response.status);
    }

    const data = await response.json();

    if (data) {
      return data;
    } else {
      alert("Response was empty, but without error");
      return null;
    }
  } catch (error) {
    console.error("Error:", error);
    return null;
  }
}

function updateStatistics() {
  getStatsFromBackEnd().then((data) => {
    if (!!data) {
      lastUpdate.textContent = new Date(
        data.LastUpdate * 1000
      ).toLocaleString();
      activePRComments.textContent = data.NumberOfActivePrComments;
      activePRTasks.textContent = data.NumberOfActivePrTasks;
    }
  });
}

async function getStatsFromBackEnd() {
  try {
    const response = await fetch("/stats", {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error("Network response was not ok: " + response.status);
    }

    const data = await response.json();

    if (data) {
      return data;
    } else {
      alert("Response was empty, but without error");
      return null;
    }
  } catch (error) {
    console.error("Error:", error);
    return null;
  }
}
