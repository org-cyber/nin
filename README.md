Below is a professional README for your project. You can copy it into a `README.md` file in the root of your repository.

```markdown
# NIN Lookup & Digital Slip Printer

An internal web tool for cybercafés to instantly retrieve a customer’s NIN details and print a high‑quality **Digital NIN Slip** on PVC cards.  
It replaces expensive third‑party services by connecting directly to a licensed NIN verification provider.

---

## ✨ Features

- **Live NIN lookups** – enter a National Identification Number and fetch name, photo, date of birth, gender, address and more.
- **Official‑style slip design** – the printed card mirrors the look of the genuine Improved NIN Slip (Premium), complete with:
  - Green header with Nigerian coat of arms
  - Photo, bilingual labels, issue date
  - QR code encoding the NIN
  - Security watermarks and diagonal NIN numbers
- **Demo mode** – test the layout and printing without spending money.
- **Zero‑cost test environment** – uses mock data that you can toggle on/off.
- **PVC‑ready** – prints on standard CR80 (85.6 mm × 54 mm) plastic cards when you have the appropriate printer.
- **Lightweight & portable** – runs on a local machine with just Node.js.

---

## 🧰 Tech Stack

| Layer       | Technology                         |
|-------------|------------------------------------|
| Backend     | Node.js (Express)                  |
| Frontend    | HTML5, CSS3, Vanilla JavaScript   |
| QR Code     | [qrcodejs](https://github.com/davidshimjs/qrcodejs) |
| Print       | Browser `window.print()` with `@page` |

---

## 📁 Project Structure

```
nin-kiosk/
├── server.js          # Express backend (proxies NIN lookups)
├── static/
│   ├── index.html     # Main UI & print card
│   └── background.png # (optional) background image for the slip
├── .env               # Your API key (never commit this)
├── package.json
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) v18 or newer
- A funded account on [ninbvnportal.com.ng](https://ninbvnportal.com.ng) (or the provider you registered with)
- Your API key from the provider’s dashboard

### 1. Clone & install

```bash
git clone https://github.com/org-cyber/nin.git
cd nin
npm install
```

### 2. Configure your API key

Create a `.env` file in the project root:

```env
CHECKMYNIN_API_KEY=your_real_api_key_here
```

> **Important:** Never commit or share this file.

### 3. Start the server

```bash
node server.js
```

The app will be available at `http://localhost:8080`.

---

## 📖 How to Use

1. Open `http://localhost:8080` in a browser.
2. You’ll see a control strip at the top:
   - **Demo Mode** (checked by default) – uses sample data (no charges).
   - **NIN input** – type the 11‑digit NIN.
   - **Lookup** button – triggers the verification.
3. When **Demo Mode** is **unchecked**, a real lookup is performed (costs ₦150 per successful verification).  
   The card instantly updates with the person’s photo, name, date of birth, gender, and issue date.
4. Click **Print Digital Slip** to open the browser’s print dialog.  
   The printed output will be a single high‑quality card (with background and borders, if the browser supports `print-color-adjust`).  
   For best results, use **Chrome** or **Edge**, and enable **Background graphics** in the print settings.
5. (Optional) If you have a PVC ID card printer, load blank CR80 cards and print directly onto them.

---

## 🔒 Security & Compliance

- **Data not stored** – the app does **not** save any personal data. Slip images are generated client‑side and discarded after printing.
- **Consent required** – every lookup sends `"consent": true`, as required by the provider and Nigerian data protection laws.
- **API key kept on the server** – the key is never exposed to the browser.
- **Rate‑limiting** – the provider enforces a 1‑minute cooldown for repeated lookups of the same NIN.

---

## 🛠️ Customisation

- **Background image** – place a file named `background.png` in the `static/` folder to add a background to the card.
- **Second image** – an extra image can be added below the card (see the note at the bottom of `index.html`).
- **Print margins** – adjust the `@page { margin: 6mm; }` value in the `<style>` block to change the whitespace around the card.

---

## 🧑‍💻 Development & Testing

- Use **Demo Mode** for layout changes – it never calls the real API.
- The server returns detailed errors if the wallet is empty, the API key is invalid, or the NIN does not exist.
- To see the raw provider response, check the server console (it logs every response).

---

## 📦 Dependencies

- [express](https://www.npmjs.com/package/express) – web framework
- [dotenv](https://www.npmjs.com/package/dotenv) – loads environment variables from `.env`
- [qrcodejs](https://github.com/davidshimjs/qrcodejs) – client‑side QR code generation

All frontend assets (fonts, QR library) are loaded from CDN.

---

## 📄 License

This project is for internal use. You may modify and distribute it as needed.

---

## 🤝 Acknowledgements

- The card design is inspired by the official Nigerian **Digital NIN Slip**.
- QR code generation by [qrcodejs](https://github.com/davidshimjs/qrcodejs).
- Background image and coat of arms illustration (if used) are based on publicly available NIMC materials.

---

**Made for cybercafé teams – save time and money on NIN printing.**
```

