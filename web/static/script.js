// web/static/script.js

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

  // Save Configuration
  saveButton.addEventListener("click", function () {
    const frequencyValue = monitoringFrequency.value;
    const commentsChecked = notifyComments.checked;
    const tasksChecked = notifyTasks.checked;
    const statusChecked = notifyStatus.checked;

    // Send the configuration to your Go backend for saving

    // For now, let's log the values to the console
    console.log("Frequency:", frequencyValue);
    console.log("Notify on Comments:", commentsChecked);
    console.log("Notify on Tasks:", tasksChecked);
    console.log("Notify on Status Changes:", statusChecked);
  });

  // Simulate Updating Statistics
  function updateStatistics() {
    // In a real application, fetch statistics from your Go backend
    // For now, let's simulate random statistics
    const randomComments = Math.floor(Math.random() * 10);
    const randomTasks = Math.floor(Math.random() * 5);

    lastUpdate.textContent = new Date().toLocaleString();
    activePRComments.textContent = randomComments;
    activePRTasks.textContent = randomTasks;
  }

  // Update statistics initially
  updateStatistics();

  // Periodically update statistics (every 30 seconds in this example)
  setInterval(updateStatistics, 30000);
});
