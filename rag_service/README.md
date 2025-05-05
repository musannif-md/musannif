# RAG Service

This service implements a Retrieval-Augmented Generation (RAG) system using Langchain and Google's Gemini API.

## Setup

1. Create a virtual environment and activate it:
```bash
conda create -n rag-env python=3.9
conda activate rag-env
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

3. Create a `.env` file with your configuration:
```bash
cp .env.example .env
```

4. Add your Gemini API key to the `.env` file:
```
GEMINI_API_KEY=your_gemini_api_key_here
```

## Usage

1. Start the API server:
```bash
python api.py
```

2. The API will be available at `http://localhost:8000`

### API Endpoints

#### Add Documents
```bash
curl -X POST "http://localhost:8000/documents" \
     -H "Content-Type: application/json" \
     -d '{"documents": ["Your document text here"]}'
```

#### Query
```bash
curl -X POST "http://localhost:8000/query" \
     -H "Content-Type: application/json" \
     -d '{"question": "Your question here"}'
```

## Features

- Document ingestion and chunking
- Vector storage using ChromaDB
- Semantic search using HuggingFace embeddings
- Question answering using Google's Gemini model
- RESTful API interface 