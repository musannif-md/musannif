from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
from rag_service import RAGService

app = FastAPI(title="RAG Service API")
rag_service = RAGService()

class DocumentRequest(BaseModel):
    documents: List[str]

class QueryRequest(BaseModel):
    question: str

@app.post("/documents")
async def add_documents(request: DocumentRequest):
    """Add documents to the vector store."""
    try:
        rag_service.create_vector_store(request.documents)
        return {"message": "Documents added successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/query")
async def query(request: QueryRequest):
    """Query the RAG system."""
    try:
        response = rag_service.query(request.question)
        return {"response": response}
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000) 