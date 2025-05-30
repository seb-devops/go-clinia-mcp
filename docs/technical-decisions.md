# Technical Decisions: LangChain Query Decomposition Tool

This document outlines the key technical decisions made during the implementation of the LangChain query decomposition tool for the MCP server.

## Overview

The goal was to integrate LangChain Go with an existing MCP (Model Context Protocol) server to provide query decomposition functionality that breaks complex user queries into smaller, actionable subqueries.

## Key Technical Decisions

### 1. LangChain Go Library Selection

**Decision**: Use `github.com/tmc/langchaingo` as the LangChain implementation.

**Reasoning**:
- Official Go port of LangChain with active maintenance
- Comprehensive agent/executor functionality 
- Good integration with OpenAI and other LLM providers
- Well-documented API and examples available
- Version 0.1.13 provides stable agent functionality

**Alternatives Considered**:
- Building custom LLM integration from scratch (rejected due to complexity)
- Using Python LangChain via subprocess calls (rejected due to performance/complexity)

### 2. Agent Architecture Choice

**Decision**: Use `OneShotZeroAgent` for query decomposition.

**Reasoning**:
- Designed for single-turn interactions without iterative tool calling
- Perfect fit for query decomposition which doesn't require multi-step reasoning
- Simpler than conversational agents for our use case
- Good balance of functionality and performance

**Alternatives Considered**:
- `ConversationalAgent` (rejected as overkill for single-turn decomposition)
- `OpenAIFunctionsAgent` (rejected as we don't need function calling)
- Custom agent implementation (rejected due to development time)

### 3. LLM Provider Strategy

**Decision**: Use OpenAI as the primary LLM provider through LangChain.

**Reasoning**:
- Existing project already uses OpenAI for embeddings
- Consistent API key management
- High-quality text generation for query analysis
- Well-supported by LangChain Go
- Good performance for natural language decomposition tasks

**Future Considerations**:
- Could be extended to support multiple providers (Anthropic, Google, etc.)
- Provider abstraction layer could be added later

### 4. Package Architecture

**Decision**: Create a separate `mcp/langchain` package for LangChain integration.

**Reasoning**:
- Follows existing project structure (`mcp/openai`, `mcp/supabase`)
- Provides clear separation of concerns
- Makes testing and maintenance easier
- Allows for future expansion of LangChain functionality
- Maintains modularity and reusability

**Structure**:
```
mcp/langchain/
├── decomposer.go      # Main implementation
├── decomposer_test.go # Test suite
└── README.md          # Package documentation
```

### 5. Response Format Design

**Decision**: Use structured JSON response with detailed metadata.

**Response Schema**:
```json
{
  "original_query": "string",
  "sub_queries": [
    {
      "id": "int",
      "query": "string", 
      "description": "string",
      "priority": "int"
    }
  ],
  "strategy": "string"
}
```

**Reasoning**:
- Provides rich metadata for each subquery
- Enables priority-based processing
- Includes strategy information for transparency
- JSON format integrates well with MCP protocol
- Extensible structure for future enhancements

### 6. Prompt Engineering Strategy

**Decision**: Use specialized prompt with clear formatting instructions.

**Prompt Structure**:
- Clear role definition ("query decomposition expert")
- Specific output format requirements (numbered list with descriptions)
- Explicit constraints (3-5 subqueries, independence requirement)
- Quality guidelines (specific, actionable, prioritized)

**Reasoning**:
- Ensures consistent output format for reliable parsing
- Provides clear guidance to the LLM
- Minimizes parsing errors and edge cases
- Balances comprehensiveness with manageability (3-5 subqueries)

### 7. Response Parsing Approach

**Decision**: Parse numbered list format with fallback handling.

**Implementation**:
- Primary: Parse "1. [Description] - [Query]" format
- Fallback: Parse "1. [Query]" format with auto-generated descriptions
- Error handling: Validate minimum subquery count

**Reasoning**:
- Reliable parsing of LLM-generated content
- Handles variations in LLM output format
- Provides graceful degradation for edge cases
- Maintains structured output even with format variations

### 8. Testing Strategy

**Decision**: Conditional tests with API key detection.

**Approach**:
- Skip integration tests when `OPENAI_API_KEY` not available
- Unit tests for parsing logic don't require API access
- Comprehensive test coverage for core functionality
- Realistic test queries of varying complexity

**Reasoning**:
- Enables testing in CI/CD environments without API keys
- Provides local testing capabilities for developers
- Maintains test reliability and consistency
- Separates unit tests from integration tests

### 9. Error Handling Strategy

**Decision**: Multi-layer error handling with context preservation.

**Layers**:
1. Input validation (empty queries, parameter checking)
2. LangChain agent errors (API failures, model errors)
3. Response parsing errors (format validation, content extraction)
4. MCP integration errors (JSON marshaling, result formatting)

**Reasoning**:
- Provides clear error messages for debugging
- Maintains system stability under various failure conditions
- Enables proper error reporting through MCP protocol
- Facilitates troubleshooting and monitoring

### 10. Integration Pattern

**Decision**: Follow existing MCP tool registration pattern.

**Implementation**:
- Tool definition using `mcp.NewTool()` 
- Handler function following MCP signature
- Registration in `setupServerAndTools()`
- Consistent error handling and logging

**Reasoning**:
- Maintains consistency with existing codebase
- Leverages established patterns and conventions
- Simplifies maintenance and future tool additions
- Provides familiar developer experience

## Performance Considerations

### Latency
- Single API call to OpenAI per decomposition request
- Typical response time: 1-3 seconds depending on query complexity
- No caching implemented initially (future enhancement opportunity)

### Resource Usage
- Minimal memory footprint (agent reused across requests)
- Network-bound operation (depends on OpenAI API availability)
- CPU usage minimal except during response parsing

### Scalability
- Stateless design enables horizontal scaling
- Rate limiting handled by OpenAI API limits
- Could be enhanced with connection pooling and caching

## Security Considerations

### API Key Management
- OpenAI API key required via environment variable
- No API key storage in code or logs
- Follows existing project security patterns

### Input Validation
- Query length and content validation
- Prevents injection attacks through input sanitization
- Error messages don't expose sensitive information

### Response Security
- No persistent storage of queries or responses
- Responses contain only decomposed query information
- No external service calls beyond OpenAI API

## Future Enhancement Opportunities

1. **Multi-Provider Support**: Add support for Anthropic, Google, etc.
2. **Caching Layer**: Implement query result caching for common patterns
3. **Custom Prompts**: Allow configurable decomposition strategies
4. **Async Processing**: Support for background query processing
5. **Metrics Collection**: Add monitoring and performance metrics
6. **Multi-Language**: Support for non-English query decomposition

## Lessons Learned

1. **API Integration**: LangChain Go provides good abstractions but requires understanding of underlying patterns
2. **Response Parsing**: LLM outputs require robust parsing with multiple fallback strategies
3. **Testing**: Conditional testing strategies essential for external API dependencies
4. **Documentation**: Comprehensive documentation crucial for complex integrations
5. **Error Handling**: Multiple error handling layers needed for reliable operation

## Conclusion

The implementation successfully integrates LangChain Go with the MCP server while maintaining code quality, testability, and extensibility. The technical decisions prioritize reliability, maintainability, and consistency with existing project patterns.

---

*Document Version: 1.0*  
*Last Updated: December 2024*  
*Author: Implementation Team* 