import os
from typing import List, Optional
from dotenv import load_dotenv
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.embeddings import HuggingFaceEmbeddings
from langchain.vectorstores import Chroma
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.schema import Document
from langchain.chains import RetrievalQA

# Load environment variables
load_dotenv()

class RAGService:
    def __init__(self):
        self.embeddings = HuggingFaceEmbeddings(
            model_name=os.getenv("EMBEDDING_MODEL", "all-MiniLM-L6-v2")
        )
        self.llm = ChatGoogleGenerativeAI(
            model=os.getenv("MODEL_NAME", "gemini-pro"),
            google_api_key=os.getenv("GEMINI_API_KEY"),
            temperature=0.7
        )
        self.vector_store = None
        self.qa_chain = None

    def create_vector_store(self, documents: List[str]):
        """Create a vector store from the provided documents."""
        # Split documents into chunks
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200
        )
        texts = text_splitter.create_documents(documents)
        
        # Create and persist vector store
        self.vector_store = Chroma.from_documents(
            documents=texts,
            embedding=self.embeddings,
            persist_directory=os.getenv("CHROMA_PERSIST_DIRECTORY", "./chroma_db")
        )
        
        # Create QA chain
        self.qa_chain = RetrievalQA.from_chain_type(
            llm=self.llm,
            chain_type="stuff",
            retriever=self.vector_store.as_retriever()
        )

    def query(self, question: str) -> str:
        """Query the RAG system with a question."""
        if not self.qa_chain:
            raise ValueError("Vector store not initialized. Call create_vector_store first.")
        
        response = self.qa_chain.run(question)
        return response

    def add_documents(self, documents: List[str]):
        """Add new documents to the existing vector store."""
        if not self.vector_store:
            self.create_vector_store(documents)
            return
        
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200
        )
        texts = text_splitter.create_documents(documents)
        
        self.vector_store.add_documents(texts)
        self.vector_store.persist() 