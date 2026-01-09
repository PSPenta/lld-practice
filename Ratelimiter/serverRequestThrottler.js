const fn3 = async () => {
  return new Promise((resolve) => {
    console.log("Starting task");
    setTimeout(
      async () => {
        console.log("Completed task");
        resolve();
      },
      Math.random() * 1000 * 2,
    );
  });
};
const fnInput = [fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3, fn3];
// Starting task
// Starting task
// Starting task
// Completed task
// Starting task
// Completed task
// Completed task
// Completed task





const limitConcurrency = async (tasks, concurrency = 3) => {
  let index = 0;

  const worker = async () => {
    while (index < tasks.length) {
      const currentTask = tasks[index++];
      await currentTask();  // Execute and wait
    }
  };

  // Create "concurrency" number of workers
  const workers = Array.from({ length: concurrency }, worker);
console.log(workers);
  // Wait for all workers to finish
  await Promise.all(workers);
};

limitConcurrency(fnInput, 3);
