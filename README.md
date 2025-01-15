# Fake News Detection

Fake News Detection is a lightweight, browser-based chatbot designed for analyzing news articles and answering general questions. It can process both plain text and links to news articles to provide insightful feedback on their credibility and summarize their content. This chatbot offers a modern dark mode interface with sleek design and intuitive functionality.

## Features

- **Analyze News Articles**:
  - Paste a full news article or provide a link to analyze and evaluate its credibility.
- **Predefined Responses**:
  - Responds to common questions with predefined answers.
- **Dark Mode Design**:
  - Sleek, professional interface optimized for readability and modern aesthetics.
- **Hybrid Input Handling**:
  - Automatically detects if the input is a URL or plain text and processes accordingly.

## Installation

To run Fake News Detection locally:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fake-news-detection.git
   ```
2. Navigate to the project directory:
   ```bash
   cd fake-news-detection
   ```
3. Open the `index.html` file in any modern web browser.

## Usage

1. Open `index.html` in your browser.
2. Type a question or paste a news article/link into the input box.
3. Click the **Send** button or press **Enter** to receive a response.

## How It Works

- **Predefined Logic**:
  - Matches specific questions (e.g., "Who made you?") to predefined answers.
- **News Analysis**:
  - Detects URLs or plain text and processes them to provide a summary and a verdict (e.g., whether the news is likely true or fake).
- **Stylish Interface**:
  - Chat messages from the user and bot are displayed in distinct, visually appealing bubbles.

## Dependencies

Fake News Detection uses the following:

- [Tailwind CSS](https://tailwindcss.com/) for responsive styling.
- [Marked.js](https://github.com/markedjs/marked) for formatting bot responses.

## Customization

1. **Modify Predefined Responses**:
   - Edit the `predefinedResponses` object in the JavaScript section of `index.html` to add or change responses.

2. **Adjust Appearance**:
   - Update the CSS styles in the `<style>` section of `index.html` to customize the chat interface.

3. **API Integration**:
   - Replace the placeholder `YOUR_API_KEY_HERE` in the `analyzeNews` function with a valid API key for advanced text analysis.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add a new feature"
   ```
4. Push to the branch:
   ```bash
   git push origin feature-name
   ```
5. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the creators of Tailwind CSS and Marked.js for their fantastic tools.
- Inspired by modern chatbot designs and functionality.

---

Feel free to reach out with questions, suggestions, or feature requests by opening an issue on GitHub!
