// Telegram WebApp
const tg = window.Telegram && window.Telegram.WebApp ? window.Telegram.WebApp : null;
let initData = "";

if (tg) {
  tg.expand();
  initData = tg.initData || "";
} else {
  console.warn("Telegram WebApp not detected. Running in dev mode.");
  // For local dev you can paste a saved initData value here if needed.
}

// --- Game logic ---

const boardElement = document.getElementById("board");
const statusText = document.getElementById("statusText");
const resultSection = document.getElementById("resultSection");
const resultTitle = document.getElementById("resultTitle");
const promoBlock = document.getElementById("promoBlock");
const promoCodeText = document.getElementById("promoCodeText");
const copyPromoBtn = document.getElementById("copyPromoBtn");
const playAgainBtn = document.getElementById("playAgainBtn");

let board = Array(9).fill(null); // "X", "O" or null
let isGameOver = false;
const PLAYER = "X";
const COMPUTER = "O";

function resetGame() {
  board = Array(9).fill(null);
  isGameOver = false;

  const cells = boardElement.querySelectorAll(".cell");
  cells.forEach((cell) => {
    cell.textContent = "";
    cell.classList.remove("x", "o", "disabled");
  });

  resultSection.hidden = true;
  promoBlock.hidden = true;
  promoCodeText.textContent = "";
  statusText.textContent = "Твой ход";
}

function checkWinner(b) {
  const winningCombinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];

  for (const [a, c, d] of winningCombinations) {
    if (b[a] && b[a] === b[c] && b[a] === b[d]) {
      return b[a];
    }
  }

  if (b.every((cell) => cell !== null)) {
    return "draw";
  }

  return null;
}

function computerMove() {
  if (isGameOver) return;

  // Small delay to feel like the computer is thinking.
  setTimeout(() => {
    let move = findBestMove(COMPUTER);
    if (move === -1) {
      move = findBestMove(PLAYER);
    }
    if (move === -1 && board[4] === null) {
      move = 4;
    }
    if (move === -1) {
      const available = board
        .map((val, idx) => (val === null ? idx : -1))
        .filter((idx) => idx !== -1);
      if (available.length > 0) {
        const randIndex = Math.floor(Math.random() * available.length);
        move = available[randIndex];
      }
    }

    if (move >= 0 && board[move] === null) {
      board[move] = COMPUTER;
      const cell = boardElement.querySelector(`.cell[data-index="${move}"]`);
      cell.textContent = "O";
      cell.classList.add("o", "disabled");
    }

    const result = checkWinner(board);
    if (result === COMPUTER) {
      isGameOver = true;
      handleGameEnd("lose");
    } else if (result === "draw") {
      isGameOver = true;
      handleGameEnd("draw");
    } else {
      statusText.textContent = "Твой ход";
    }
  }, 500);
}

function findBestMove(playerSymbol) {
  const b = [...board];
  const winningCombinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];

  for (const [a, c, d] of winningCombinations) {
    const line = [b[a], b[c], b[d]];
    const marks = line.filter((v) => v === playerSymbol).length;
    const empties = line.filter((v) => v === null).length;

    if (marks === 2 && empties === 1) {
      if (b[a] === null) return a;
      if (b[c] === null) return c;
      if (b[d] === null) return d;
    }
  }

  return -1;
}

function handleCellClick(e) {
  const cell = e.target;
  if (!cell.classList.contains("cell")) return;
  const index = parseInt(cell.getAttribute("data-index"), 10);

  if (isGameOver || board[index] !== null) return;

  board[index] = PLAYER;
  cell.textContent = "X";
  cell.classList.add("x", "disabled");

  const result = checkWinner(board);
  if (result === PLAYER) {
    isGameOver = true;
    handleGameEnd("win");
  } else if (result === "draw") {
    isGameOver = true;
    handleGameEnd("draw");
  } else {
    statusText.textContent = "Ход компьютера…";
    computerMove();
  }
}

boardElement.addEventListener("click", handleCellClick);

function handleGameEnd(result) {
  let text = "";
  let promoCode = "";

  switch (result) {
    case "win":
      text = "Ты выиграла! 🎉";
      promoCode = generatePromoCode();
      promoCodeText.textContent = promoCode;
      promoBlock.hidden = false;
      break;
    case "lose":
      text = "В этот раз выиграл компьютер. Попробуешь ещё раз?";
      promoBlock.hidden = true;
      break;
    case "draw":
      text = "Ничья, давай ещё раз?";
      promoBlock.hidden = true;
      break;
  }

  statusText.textContent = "";
  resultTitle.textContent = text;
  resultSection.hidden = false;

  sendResult(result, promoCode);
}

function generatePromoCode() {
  const n = Math.floor(Math.random() * 100000);
  return n.toString().padStart(5, "0");
}

async function sendResult(result, promoCode) {
  try {
    const payload = {
      result,
      promoCode,
      initData,
    };

    const resp = await fetch("/api/result", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!resp.ok) {
      console.error("Failed to send result:", await resp.text());
    } else {
      const data = await resp.json();
      console.log("Result sent:", data);
    }
  } catch (err) {
    console.error("Error sending result:", err);
  }
}

playAgainBtn.addEventListener("click", () => {
  resetGame();
});

copyPromoBtn.addEventListener("click", async () => {
  const code = promoCodeText.textContent.trim();
  if (!code) return;

  try {
    await navigator.clipboard.writeText(code);
    if (tg && tg.showPopup) {
      tg.showPopup({
        title: "Скопировано",
        message: "Промокод скопирован в буфер обмена 💕",
        buttons: [{ type: "close" }],
      });
    } else {
      alert("Промокод скопирован: " + code);
    }
  } catch (err) {
    console.error("Clipboard error:", err);
  }
});

// Initialize
resetGame();
