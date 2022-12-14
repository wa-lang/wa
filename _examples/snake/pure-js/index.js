const domID = 'wa-canvas';
const canvas = document.getElementById(domID);
const ctx = canvas.getContext('2d');

// canvas 宽高和格子大小
const width = 500;
const height = 500;
const gridSize = 20;

// 蛇的位置、方向和食物的位置
const gameState = {
  snake: [{ x: 0, y: 0 }],
  dir: 'right',
  food: {
    x: Math.floor(Math.random() * (width / gridSize)),
    y: Math.floor(Math.random() * (height / gridSize))
  }
};

// 方向键对应的方向
const KEYS = {
  'ArrowUp': 'up',
  'ArrowDown': 'down',
  'ArrowLeft': 'left',
  'ArrowRight': 'right'
}

// 设置 canvas 的宽高
canvas.width = width;
canvas.height = height;

// 绘制蛇
function drawSnake() {
  // 循环遍历贪吃蛇的每一节身体
  gameState.snake.forEach((segment, index) => {
    let x = segment.x * gridSize;
    let y = segment.y * gridSize;
    const imageData = ctx.createImageData(gridSize, gridSize);
    for (let i = 0; i < imageData.data.length; i += 4) {
      imageData.data[i] = 0;
      imageData.data[i + 1] = 0;
      imageData.data[i + 2] = 0;
      imageData.data[i + 3] = 255;
    }
    ctx.putImageData(imageData, x, y);
  });
}

// 绘制食物
function drawFood() {
  const x = gameState.food.x * gridSize;
  const y = gameState.food.y * gridSize;
  const imageData = ctx.createImageData(gridSize, gridSize);
  for (let i = 0; i < imageData.data.length; i += 4) {
    imageData.data[i] = 32;
    imageData.data[i + 1] = 178;
    imageData.data[i + 2] = 179;
    imageData.data[i + 3] = 255;
  }
  ctx.putImageData(imageData, x, y);
}


// 更新蛇的位置和食物的位置
function updateGame() {
  // 蛇头位置
  const head = gameState.snake[0];

  // 根据方向移动蛇头的位置
  if (gameState.dir === 'right') {
    head.x += 1;
  } else if (gameState.dir === 'left') {
    head.x -= 1;
  } else if (gameState.dir === 'up') {
    head.y -= 1;
  } else if (gameState.dir === 'down') {
    head.y += 1;
  }

  // 吃到食物就增加一节身体并生成新的食物
  if (head.x === gameState.food.x && head.y === gameState.food.y) {
    gameState.snake.push({ x: head.x, y: head.y });
    gameState.food = generateFoodLocation();
  }

  // 移动蛇身体
  setTimeout(() => {
    for (let i = gameState.snake.length - 1; i > 0; i--) {
      gameState.snake[i].x = gameState.snake[i - 1].x;
      gameState.snake[i].y = gameState.snake[i - 1].y;
    }
  }, 0);

  // 如果蛇头撞墙游戏结束
  if (
    head.x < 0 ||
    head.x > (width / gridSize) ||
    head.y < 0 ||
    head.y > (height / gridSize)
  ) {
    clearInterval(timer);
    alert('游戏结束，请刷新页面重新开始');
    return;
  }
}

// 绘制蛇和食物
function drawGame() {
  // 清空 canvas
  ctx.clearRect(0, 0, width, height);
  // 绘制贪吃蛇
  drawSnake();
  // 绘制食物
  drawFood()
}


// 更新并绘制
function gameLoop() {
  updateGame();
  drawGame();
}

// 处理键盘事件
function handleKeyDown(event) {
  // 同一方向不允许反向移动
  if (
    (gameState.dir === 'up' && event.key === 'ArrowDown') ||
    (gameState.dir === 'down' && event.key === 'ArrowUp') ||
    (gameState.dir === 'left' && event.key === 'ArrowRight') ||
    (gameState.dir === 'right' && event.key === 'ArrowLeft')
  ) { return; }

  // 更新方向
  if (KEYS[event.key]) {
    gameState.dir = KEYS[event.key];
  }
}

// 生成食物的位置，保证食物不会出现在贪吃蛇的身体上
function generateFoodLocation() {
  let foodLocation = randomLocation();
  while (gameState.snake.some(segment => segment.x === foodLocation.x && segment.y === foodLocation.y)) {
    foodLocation = randomLocation();
  }
  return foodLocation;
}

// 随机生成食物的位置
function randomLocation() {
  return {
    x: Math.floor(Math.random() * (width / gridSize)),
    y: Math.floor(Math.random() * (height / gridSize))
  };
}

// 监听键盘事件
window.addEventListener('keydown', handleKeyDown);
// 每 160 毫秒调用一次 gameLoop 函数
const timer = setInterval(gameLoop, 160);

