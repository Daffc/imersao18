import { Controller, Get, Post, Body, Patch, Param, Delete, HttpCode, UseGuards } from '@nestjs/common';
import { EventsService } from '@app/core/events/events.service';
import { CriarEventoRequest } from './request/criar-evento.request';
import { AtualizarEventoRequest } from './request/atualizar-evento.request';
import { ReservarLugarRequest } from './request/reservar-lugar.request';
import { TicketKind } from '@prisma/client';
import { AuthGuard } from '@app/core/auth/auth.guard';

@Controller('eventos')
export class EventosController {
  constructor(private readonly EventosService: EventsService) { }

  @Post()
  create(@Body() criarEventoRequest: CriarEventoRequest) {
    return this.EventosService.create({
      name: criarEventoRequest.nome,
      description: criarEventoRequest.descricao,
      date: criarEventoRequest.data,
      price: criarEventoRequest.preco,
    });
  }

  @Get()
  findAll() {
    return this.EventosService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.EventosService.findOne(id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() atualizarEventoRequest: AtualizarEventoRequest) {
    return this.EventosService.update(id, ({
      name: atualizarEventoRequest.nome,
      description: atualizarEventoRequest.descricao,
      date: atualizarEventoRequest.data,
      price: atualizarEventoRequest.preco,
    }));
  }

  @HttpCode(204)
  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.EventosService.remove(id);
  }

  @UseGuards(AuthGuard)
  @Post(':id/reservar')
  reserveSpots(@Body() reservarLugarRequest: ReservarLugarRequest, @Param(':id') eventId: string) {
    return this.EventosService.reserveSpot({
      eventId,
      spots: reservarLugarRequest.lugares,
      ticket_kind: reservarLugarRequest.tipo_ingresso === "inteira" ? TicketKind.full : TicketKind.half,
      email: reservarLugarRequest.email
    })
  }
}
